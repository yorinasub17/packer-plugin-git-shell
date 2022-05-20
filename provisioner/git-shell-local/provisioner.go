//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package git_shell_local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	shell_local_provisioner "github.com/hashicorp/packer/provisioner/shell-local"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

type Config struct {
	common.GitConfig    `mapstructure:",squash"`
	common.ScriptConfig `mapstructure:",squash"`

	ctx interpolate.Context
}

type Provisioner struct {
	config Config
}

func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.provisioner.git-shell-local",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{}) (returnErr error) {
	// Clone repo to a temporary directory that is cleaned up later
	cloneDir, err := ioutil.TempDir("", "packer-git-shell-*")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(cloneDir); err != nil {
			returnErr = multierror.Append(returnErr, err)
		}
	}()

	gitOpts := p.config.GitConfig.GetGitOptions()
	ui.Say(
		fmt.Sprintf("Cloning repo %s at ref %s to dir %s", gitOpts.RepoURL, gitOpts.Ref, cloneDir),
	)
	if err := common.CloneAndCheckout(gitOpts, cloneDir); err != nil {
		return err
	}

	// Run the script through the shell-local provisioner now that we have the script locally and thus won't be any
	// different. Since the config property of the provisioner is internal, we need to use the map[string]interface{}
	// representation of the config so that we can use Prepare to get the config representation.
	// We use inline scripts to construct the script calls so that we can call each script with its args and environment
	// variables.
	scriptCalls := []string{}
	for _, script := range p.config.ScriptConfig.Scripts {
		scriptCalls = append(
			scriptCalls,
			fmt.Sprintf(
				"%s %s %s",
				strings.Join(script.EnvironmentVars, " "),
				filepath.Join(cloneDir, script.Path),
				strings.Join(script.Args, " "),
			),
		)
	}

	shellConfig := map[string]interface{}{
		"inline": scriptCalls,
	}
	shellProvisioner := &shell_local_provisioner.Provisioner{}
	if err := shellProvisioner.Prepare(shellConfig); err != nil {
		return err
	}

	return shellProvisioner.Provision(ctx, ui, comm, generatedData)
}

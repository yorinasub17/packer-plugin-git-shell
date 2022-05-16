//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package git_shell

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
	shell_provisioner "github.com/hashicorp/packer/provisioner/shell"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

type Config struct {
	common.Config `mapstructure:",squash"`

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
		PluginType:         "packer.provisioner.git-shell",
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

func (p *Provisioner) Provision(
	ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{},
) (returnErr error) {
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

	gitOpts := p.config.Config.GetGitOptions()
	ui.Say(
		fmt.Sprintf("Cloning repo %s at ref %s to dir %s", gitOpts.RepoURL, gitOpts.Ref, cloneDir),
	)
	if err := common.CloneAndCheckout(gitOpts, cloneDir); err != nil {
		return err
	}

	// Call each script in a loop, running through the shell provisioner. This is less efficient than passing the shell
	// provisioner all the scripts, but since the shell provisioner doesn't have a concept of script args, we have to
	// cheat with the execute_command field, which apply to all scripts passed to the shell provisioner. As such, we
	// need to switch the context for each script, requiring a separate provisioner call.
	for _, script := range p.config.Config.Scripts {
		if err := runScript(ctx, ui, comm, generatedData, cloneDir, script); err != nil {
			return err
		}
	}
	return nil
}

func runScript(
	ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{},
	cloneDir string, script common.Script,
) error {
	// Since the config property of the provisioner is internal, we need to use the map[string]interface{}
	// representation of the config so that we can use Prepare to get the config representation.
	execCmd := fmt.Sprintf("chmod +x {{.Path}}; {{.Vars}} {{.Path}} %s",
		strings.Join(script.Args, " "),
	)
	shellConfig := map[string]interface{}{
		"script":           filepath.Join(cloneDir, script.Path),
		"environment_vars": script.EnvironmentVars,
		"execute_command":  execCmd,
	}
	shellProvisioner := &shell_provisioner.Provisioner{}
	if err := shellProvisioner.Prepare(shellConfig); err != nil {
		return err
	}

	// Run the script through the shell provisioner now that we have the script locally and thus won't be any different.
	return shellProvisioner.Provision(ctx, ui, comm, generatedData)
}

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package git_shell_file

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	file_provisioner "github.com/hashicorp/packer/provisioner/file"

	"github.com/yorinasub17/packer-plugin-git-shell/provisioner/common"
)

type Config struct {
	common.GitConfig `mapstructure:",squash"`

	// Files is a list of blocks that specify which files from the repo should be uploaded, and to where. The
	// files will be uploaded in the order in which the blocks are defined. At least one file block must be defined.
	Files []common.File `mapstructure:"file"`

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
		PluginType:         "packer.provisioner.git-shell-file",
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
	cloneDir, err := ioutil.TempDir("", "packer-git-shell-file-*")
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

	// Upload each file in a loop, running through the file provisioner. We do this because the file provisioner can
	// only upload to a single destination, and it is not guaranteed all the files are going to the same place.
	for _, file := range p.config.Files {
		if err := uploadFile(ctx, ui, comm, generatedData, cloneDir, file); err != nil {
			return err
		}
	}
	return nil
}

func uploadFile(
	ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{},
	cloneDir string, file common.File,
) error {
	// Since the config property of the provisioner is internal, we need to use the map[string]interface{}
	// representation of the config so that we can use Prepare to get the config representation.
	fileConfig := map[string]interface{}{
		"source":      filepath.Join(cloneDir, file.Path),
		"destination": file.Destination,
	}
	fileProvisioner := &file_provisioner.Provisioner{}
	if err := fileProvisioner.Prepare(fileConfig); err != nil {
		return err
	}

	// Run the file through the file provisioner now that we have the file locally and thus won't be any different.
	return fileProvisioner.Provision(ctx, ui, comm, generatedData)
}

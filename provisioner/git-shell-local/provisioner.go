//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package git_shell_local

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	// Source is the git URL (e.g., https://github.com/yorinasub17/packer-git-shell-example.git) where the scripts are
	// located.
	Source string `mapstructure:"source"`
	// Ref is the git ref to checkout when sourcing the scripts.
	Ref string `mapstructure:"ref"`
	// Script is the relative path in the git repo where the script to run is located.
	Script string `mapstructure:"script"`
	// Args is the script args to pass when executing the script.
	Args []string `mapstructure:"args"`

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

func (p *Provisioner) Provision(_ context.Context, ui packer.Ui, _ packer.Communicator, generatedData map[string]interface{}) error {
	ui.Say(
		fmt.Sprintf(
			"local provisioner args: %s//%s?ref=%s %v",
			p.config.Source, p.config.Script, p.config.Ref, p.config.Args,
		),
	)
	return nil
}

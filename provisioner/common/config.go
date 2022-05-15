package common

import (
	packer_common "github.com/hashicorp/packer-plugin-sdk/common"
)

type Config struct {
	packer_common.PackerConfig `mapstructure:",squash"`

	// Source is the git URL (e.g., https://github.com/yorinasub17/packer-git-shell-example.git) where the scripts are
	// located.
	Source string `mapstructure:"source"`
	// Ref is the git ref to checkout when sourcing the scripts.
	Ref string `mapstructure:"ref"`
	// Script is the relative path in the git repo where the script to run is located.
	Script string `mapstructure:"script"`
	// Args is the script args to pass when executing the script.
	Args []string `mapstructure:"args"`
}

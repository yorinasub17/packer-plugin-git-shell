package common

import (
	packer_common "github.com/hashicorp/packer-plugin-sdk/common"
)

type GitConfig struct {
	packer_common.PackerConfig `mapstructure:",squash"`

	// Source is the git URL (e.g., https://github.com/yorinasub17/packer-git-shell-example.git) where the scripts are
	// located.
	Source string `mapstructure:"source" required:"true"`
	// Ref is the git ref to checkout when sourcing the scripts.
	Ref string `mapstructure:"ref" required:"true"`
	// UsernameEnvVar is the name of the environment variable to lookup for the username to use when authing to the git
	// repo. Defaults to GIT_USERNAME.
	UsernameEnvVar string `mapstructure:"username_env_var"`
	// PasswordEnvVar is the name of the environment variable to lookup for the password to use when authing to the git
	// repo. If unset, defaults to GIT_PASSWORD.
	PasswordEnvVar string `mapstructure:"password_env_var"`
}

type ScriptConfig struct {
	// Scripts is a list of blocks that specify which scripts from the repo should be called, and with what args. The
	// scripts will be called in the order in which the blocks are defined. At least one script block must be defined.
	Scripts []Script `mapstructure:"script"`
}

func (cfg GitConfig) GetGitOptions() GitOptions {
	return GitOptions{
		RepoURL:        cfg.Source,
		Ref:            cfg.Ref,
		UsernameEnvVar: cfg.UsernameEnvVar,
		PasswordEnvVar: cfg.PasswordEnvVar,
	}
}

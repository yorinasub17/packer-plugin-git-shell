package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"

	git_shell "github.com/yorinasub17/packer-plugin-git-shell/provisioner/git-shell"
	git_shell_local "github.com/yorinasub17/packer-plugin-git-shell/provisioner/git-shell-local"
)

var (
	// Version is the main version number that is being run at the moment. This is automatically updated by goreleaser
	// at build time.
	Version = "1.0.0"

	// VersionPrerelease is a pre-release marker for the Version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"
)

var PluginVersion *version.PluginVersion

func init() {
	PluginVersion = version.InitializePluginVersion(Version, VersionPrerelease)
}

func main() {
	pps := plugin.NewSet()
	pps.RegisterProvisioner(plugin.DEFAULT_NAME, new(git_shell.Provisioner))
	pps.RegisterProvisioner("local", new(git_shell_local.Provisioner))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

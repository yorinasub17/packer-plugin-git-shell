//go:generate packer-sdc mapstructure-to-hcl2 -type Script

package common

type Script struct {
	// Path is the relative path in the git repo where the script to run is located.
	Path string `mapstructure:"path"`
	// Args is the script args to pass when executing the script.
	Args []string `mapstructure:"args"`
	// EnvironmentVars is the list of environment variables that should be set when the script is called. Each entry is
	// a key=value string, setting the environment variable `key` to the given `value`.
	EnvironmentVars []string `mapstructure:"environment_vars"`
}
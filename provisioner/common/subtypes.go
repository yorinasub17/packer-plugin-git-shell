//go:generate packer-sdc mapstructure-to-hcl2 -type Script,File

package common

type Script struct {
	// Path is the relative path in the git repo where the script to run is located.
	Path string `mapstructure:"path" required:"true"`
	// Args is the script args to pass when executing the script.
	Args []string `mapstructure:"args"`
	// EnvironmentVars is the list of environment variables that should be set when the script is called. Each entry is
	// a key=value string, setting the environment variable `key` to the given `value`.
	EnvironmentVars []string `mapstructure:"environment_vars"`
}

type File struct {
	// Path is the relative path in the git repo where the file or folder to upload is located.
	Path string `mapstructure:"path" required:"true"`
	// Destination is the path on the machine where the file or folder will be uploaded to.
	Destination string `mapstructure:"destination" required:"true"`
}

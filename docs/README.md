# Git Shell Plugins

<!--
  Include a short overview about the plugin.

  This document is a great location for creating a table of contents for each
  of the components the plugin may provide. This document should load automatically
  when navigating to the docs directory for a plugin.

-->

## Installation

### Using pre-built releases

#### Using the `packer init` command

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    name = {
      version = ">= 0.0.1"
      source  = "github.com/yorinasub17/git-shell"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/yorinasub17/packer-plugin-git-shell/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `make` from the root
directory. Upon successful compilation, a `packer-plugin-git-shell` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


## Plugin Contents

### Provisioners

- [git-shell](/docs/provisioners/provisioner-git-shell.mdx) - The git-shell provisioner can be used to remotely fetch a provisioner script from git.
- [git-shell-local](/docs/provisioners/provisioner-git-shell-local.mdx) - The git-shell-local provisioner can be used to remotely fetch a provisioner script from git.

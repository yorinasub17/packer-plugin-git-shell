# Git Shell Packer Plugin

The `git-shell` plugin can be used with HashiCorp [Packer](https://www.packer.io/) to provision machines using shell
scripts stored in git repositories. For the full list of available features for this plugin see [docs](/docs).

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    git-shell = {
      version = ">= 1.0.1"
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


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-git-shell` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### Configuration

For more information on how to configure the plugin, please read the
documentation located in the [`docs/`](docs) directory.


## Contributing

* If you think you've found a bug in the code or you have a question regarding the usage of this software, please reach
  out to us by [opening an issue](https://github.com/yorinasub17/packer-plugin-git-shell/issues) in this GitHub
  repository.
* Contributions to this project are welcome: if you want to add a feature or a fix a bug, please do so by [opening a
  Pull Request](https://github.com/yorinasub17/packer-plugin-git-shell/pulls) in this GitHub repository. In case of
  feature contribution, we kindly ask you to open an issue to discuss it beforehand.

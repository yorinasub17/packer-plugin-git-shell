packer {
  required_plugins {
    git-shell = {
      version = ">=v0.1.0"
      source  = "github.com/yorinasub17/git-shell"
    }
  }
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = ["sources.null.basic-example"]

  provisioner "git-shell" {
    source = "https://github.com/yorinasub17/packer-git-shell-example.git"
    ref    = "main"

    script {
      path = "scripts/echo-to-stderr"
      args = [
        "$(lsb_release -a)",
      ]
    }

    script {
      path = "scripts/echo-from-env"
      environment_vars = ["TEXT='こんにちは世界'"]
    }
  }

  provisioner "git-shell-local" {
    source = "https://github.com/yorinasub17/packer-git-shell-example.git"
    ref    = "main"

    script {
      path = "scripts/echo-to-stderr"
      args = [
        "$(lsb_release -a)",
      ]
    }

    script {
      path = "scripts/echo-from-env"
      environment_vars = ["TEXT='こんにちは世界'"]
    }
  }
}

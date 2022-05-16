packer {
  required_plugins {
    git-shell = {
      version = ">=v0.1.0"
      source  = "github.com/yorinasub17/git-shell"
    }
  }
}

# Use docker to avoid the need to open SSH locally. Otherwise, the git-shell provisioner will fail since it needs to run
# on the target machine.
source "docker" "example" {
  image  = "ubuntu"
  commit = true
}

build {
  sources = ["sources.docker.example"]

  provisioner "git-shell" {
    source = "https://github.com/yorinasub17/packer-git-shell-example.git"
    ref    = "main"

    script {
      path = "scripts/echo-to-stderr"
      args = [
        "Hello world",
      ]
    }

    script {
      path             = "scripts/echo-from-env"
      environment_vars = ["TEXT='こんにちは世界'"]
    }
  }

  provisioner "git-shell-local" {
    source = "https://github.com/yorinasub17/packer-git-shell-example.git"
    ref    = "main"

    script {
      path = "scripts/echo-to-stderr"
      args = [
        "Hello world",
      ]
    }

    script {
      path             = "scripts/echo-from-env"
      environment_vars = ["TEXT='こんにちは世界'"]
    }
  }
}

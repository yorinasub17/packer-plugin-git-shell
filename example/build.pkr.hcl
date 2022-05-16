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
        "hello",
        "world",
      ]
    }

    script {
      path = "scripts/echo-to-stderr"
      args = [
        "こんにちは",
        "世界",
      ]
    }
  }
}

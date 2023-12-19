# Use docker to avoid the need to open SSH locally. Otherwise, the git-shell provisioner will fail since it needs to run
# on the target machine.
packer {
  required_plugins {
    docker = {
      source  = "github.com/hashicorp/docker"
      version = "~> 1"
    }
  }
}

source "docker" "example" {
    image = "alpine"
    commit = true
}

build {
  sources = [
    "source.docker.example"
  ]

  provisioner "git-shell-file" {
    source           = "https://github.com/yorinasub17/packer-plugin-git-shell-automated-testing-private.git"
    ref              = "main"
    username_env_var = "TEST_GIT_USERNAME"
    password_env_var = "TEST_GIT_TOKEN"

    file {
      path        = "files/hello"
      destination = "/tmp/hello"
    }
  }

  provisioner "shell" {
    inline = [
      "cat /tmp/hello",
    ]
  }
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "git-shell-local" {
    source           = "https://github.com/yorinasub17/packer-plugin-git-shell-automated-testing-private.git"
    ref              = "main"
    username_env_var = "TEST_GIT_USERNAME"
    password_env_var = "TEST_GIT_TOKEN"

    script {
      path = "scripts/echo-private"
    }
  }
}

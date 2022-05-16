source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "git-shell-local" {
    source = "https://github.com/yorinasub17/packer-plugin-git-shell-automated-testing.git"
    ref    = "test"

    script {
      path = "scripts/echo-test"
    }

    script {
      path = "scripts/echo-public"
      args = [
        "hello",
        "world",
      ]
    }
  }
}

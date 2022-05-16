# Use docker to avoid the need to open SSH locally. Otherwise, the git-shell provisioner will fail since it needs to run
# on the target machine.
source "docker" "example" {
    image = "alpine"
    commit = true
}

build {
  sources = [
    "source.docker.example"
  ]

  provisioner "git-shell" {
    source = "https://github.com/yorinasub17/packer-plugin-git-shell-automated-testing.git"
    ref    = "main"

    script {
      path = "scripts/echo-public"
      args = [
        "hello",
        "world",
      ]
    }
  }
}

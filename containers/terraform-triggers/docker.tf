locals {
  docker_image_tag = sha256(join("", [for f in fileset(path.module, "docker/*") : filesha256(f)]))
}

resource "docker_image" "main" {
  name = "${scaleway_container_namespace.main.registry_endpoint}/server:${local.docker_image_tag}"
  build {
    context  = "docker"
    platform = "amd64"
  }
  triggers = {
    tag = local.docker_image_tag
  }
}

resource "docker_registry_image" "main" {
  name          = docker_image.main.name
  keep_remotely = true # keep old images
}

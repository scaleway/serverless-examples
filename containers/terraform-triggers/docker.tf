locals {
  docker_image_tag = sha256(join("", [for f in fileset(path.module, "docker/${var.container_language}/*") : filesha256(f)]))
}

resource "docker_image" "main" {
  name = "${scaleway_container_namespace.main.registry_endpoint}/server-${var.container_language}:${local.docker_image_tag}"
  build {
    context  = "docker/${var.container_language}"
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

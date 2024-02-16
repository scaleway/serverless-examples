resource "scaleway_registry_namespace" "main" {
  name       = "mongo-example-${random_string.suffix.result}"
  region     = "fr-par"
  project_id = var.project_id
}

resource "docker_image" "main" {
  name = "${scaleway_registry_namespace.main.endpoint}/mongo-example:${var.image_version}"
  build {
    context = "${path.cwd}/../container"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.main.name}"
  }
}

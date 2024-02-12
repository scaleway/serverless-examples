resource "scaleway_registry_namespace" "main" {
  name       = "serverless-jobs-example"
  region     = var.region
  project_id = var.project_id
}

resource "docker_image" "main" {
  name = "${scaleway_registry_namespace.main.endpoint}/jobs-hello:${var.image_version}"
  build {
    context = "${path.cwd}/../image"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.main.name}"
  }
}

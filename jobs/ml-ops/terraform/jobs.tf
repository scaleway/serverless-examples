resource "docker_image" "data_loader_image" {
  name = "${scaleway_registry_namespace.data_loader_image_registry.endpoint}/data-loader:${var.image_version}"
  build {
    context = "${path.cwd}/../s3/data-store"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.data_loader_image.name}"
  }
}

resource "docker_image" "ml_job_image" {
  name = "${scaleway_registry_namespace.ml_job_image_registry.endpoint}/ml-job:${var.image_version}"
  build {
    context = "${path.cwd}/../job"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.ml_job_image.name}"
  }
}

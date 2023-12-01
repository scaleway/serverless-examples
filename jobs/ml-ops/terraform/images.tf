resource "scaleway_registry_namespace" "main" {
  name       = "ml-ops-example-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}

resource "docker_image" "inference" {
  name = "${scaleway_registry_namespace.main.endpoint}/inference:0.0.1"
  build {
    context = "${path.cwd}/../inference"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.inference.name}"
  }
}

resource "docker_image" "data" {
  name = "${scaleway_registry_namespace.main.endpoint}/data:${var.image_version}"
  build {
    context = "${path.cwd}/../data"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.data.name}"
  }
}

resource "docker_image" "training" {
  name = "${scaleway_registry_namespace.main.endpoint}/training:${var.image_version}"
  build {
    context = "${path.cwd}/../training"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.training.name}"
  }
}

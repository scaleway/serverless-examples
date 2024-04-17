resource "scaleway_registry_namespace" "main" {
  name       = "ml-ops-example-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}

resource "docker_image" "inference" {
  name = "${scaleway_registry_namespace.main.endpoint}/inference:${var.image_version}"
  build {
    context = "${path.cwd}/../inference"
    no_cache = true
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.inference.name}"
  }
}

resource "docker_image" "data" {
  name = "${scaleway_registry_namespace.main.endpoint}/data:${var.image_version}"
  build {
    context = "${path.cwd}/../data"
    no_cache = true
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.data.name}"
  }
}

resource "docker_image" "training" {
  name = "${scaleway_registry_namespace.main.endpoint}/training:${var.image_version}"
  build {
    context = "${path.cwd}/../training"
    no_cache = true
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.training.name}"
  }
}

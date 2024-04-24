resource "scaleway_registry_namespace" "main" {
  name       = "ifr-${lower(replace(var.hf_model_file_name, "/[.]|[_]/", "-"))}-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}

resource "docker_image" "inference" {
  name = "${scaleway_registry_namespace.main.endpoint}/inference-with-huggingface:${var.image_version}"
  build {
    context = "${path.cwd}/../"
    no_cache = true
    build_args = {
     MODEL_DOWNLOAD_SOURCE : var.hf_model_download_source
    }
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.inference.name}"
  }
}

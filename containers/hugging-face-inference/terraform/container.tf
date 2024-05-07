resource "scaleway_container_namespace" "main" {
  name        = "ifr-${lower(replace(var.hf_model_file_name, "/[.]|[_]/", "-"))}-${random_string.random_suffix.result}"
  description = "Inference using Hugging Face models"
}

resource "scaleway_container" "inference-hugging-face" {
  name           = "inference"
  description    = "Inference serving API using a Hugging Face model"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = docker_image.inference.name
  environment_variables = {
    "MODEL_FILE_NAME" = var.hf_model_file_name
  }
  port           = 80
  cpu_limit      = 2240
  memory_limit   = 4096
  min_scale      = 1
  max_scale      = 1
  deploy   = true
}

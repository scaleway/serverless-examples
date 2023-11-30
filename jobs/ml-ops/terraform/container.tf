provider "docker" {
  host = "unix:///var/run/docker.sock"

  registry_auth {
    address  = scaleway_registry_namespace.inference_api_image_registry.endpoint
    username = "nologin"
    password = var.secret_key
  }

  registry_auth {
    address  = scaleway_registry_namespace.data_loader_image_registry.endpoint
    username = "nologin"
    password = var.secret_key
  }

  registry_auth {
    address  = scaleway_registry_namespace.ml_job_image_registry.endpoint
    username = "nologin"
    password = var.secret_key
  }
}

resource "docker_image" "inference_api_image" {
  name = "${scaleway_registry_namespace.inference_api_image_registry.endpoint}/inference-api:${var.image_version}"
  build {
    context = "${path.cwd}/../container/inference-api"
  }

  provisioner "local-exec" {
    command = "docker push ${docker_image.inference_api_image.name}"
  }
}

resource "scaleway_container_namespace" "inference_api_namespace" {
  name        = "ml-inference-${random_string.random_suffix.result}"
  description = "Serving inference models deployed as serverless containers"
}

resource "scaleway_container" "inference_api_container" {
  name           = "inference-api-${random_string.random_suffix.result}"
  description    = "Serving an inference API"
  namespace_id   = scaleway_container_namespace.inference_api_namespace.id
  registry_image = docker_image.inference_api_image.name
  port           = 80
  cpu_limit      = 1120
  memory_limit   = 2048
  min_scale      = 1
  max_scale      = 5
  environment_variables = {
    "MODEL_REGISTRY" = scaleway_object_bucket.model_registry.name
    "MAIN_REGION"    = var.region
    "MODEL_FILE"     = "classifier.pkl"
  }
  secret_environment_variables = {
    "SCW_ACCESS_KEY" = var.access_key
    "SCW_SECRET_KEY" = var.secret_key
  }
  privacy  = "private"
  protocol = "http1"
  deploy   = true
}

resource "scaleway_container_token" "inference_api_token" {
  container_id = scaleway_container.inference_api_container.id
}

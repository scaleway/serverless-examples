resource "scaleway_container_namespace" "inference_api_namespace" {
  name = "ml-serving"
  description = "Serving inference models deployed as serverless containers"
}

resource "scaleway_container" "inference_api_container" {
  name = "inference-api"
  description = "Serving an inference API"
  namespace_id = scaleway_container_namespace.inference_api_namespace.id
  registry_image = var.registry_image
  port = 80
  cpu_limit = 1120
  memory_limit = 2048
  min_scale = 1
  max_scale = 5
  privacy = "private"
  protocol = "http1"
  deploy = true
}
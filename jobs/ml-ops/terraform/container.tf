resource "scaleway_container_namespace" "main" {
  name        = "ml-ops-example-${random_string.random_suffix.result}"
  description = "MLOps example"
}

resource "scaleway_container" "inference" {
  name           = "inference"
  description    = "Inference serving API"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = docker_image.inference.name
  port           = 80
  cpu_limit      = 2000
  memory_limit   = 2048
  min_scale      = 1
  max_scale      = 5
  environment_variables = {
    "S3_BUCKET_NAME" = scaleway_object_bucket.main.name
    "S3_URL"         = var.s3_url
    "SCW_REGION"     = var.region
  }
  secret_environment_variables = {
    "SCW_ACCESS_KEY" = var.access_key
    "SCW_SECRET_KEY" = var.secret_key
  }
  deploy   = true
}

output "endpoint" {
  value = scaleway_container.inference.domain_name
}

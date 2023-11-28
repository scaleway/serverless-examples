resource "scaleway_registry_namespace" "inference_api_image_registry" {
  name = "inference-api-images"
  region = var.region
  project_id = var.project_id
}

resource "scaleway_registry_namespace" "ml_job_image_registry" {
  name = "ml-job-images"
  region = var.region
  project_id = var.project_id
}
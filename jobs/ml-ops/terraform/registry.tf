resource "scaleway_registry_namespace" "inference_api_image_registry" {
  name       = "inference-api-images-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}

resource "scaleway_registry_namespace" "data_loader_image_registry" {
  name       = "data-loder-images-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}

resource "scaleway_registry_namespace" "ml_job_image_registry" {
  name       = "ml-job-images-${random_string.random_suffix.result}"
  region     = var.region
  project_id = var.project_id
}


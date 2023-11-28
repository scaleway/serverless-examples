output "api_inference_registry_namespace_endpoint" {
  value = scaleway_registry_namespace.inference_api_image_registry.endpoint
}

output "ml_job_registry_namespae_endpoint" {
  value = scaleway_registry_namespace.ml_job_image_registry.endpoint
}

output "data_store_bucket_endpoint" {
  value = scaleway_object_bucket.data_store_bucket.endpoint
}

output "model_registry_bucket_endpoint" {
  value = scaleway_object_bucket.model_registry_bucket.endpoint
}

output "performance_monitoring_bucket_endpoint" {
  value = scaleway_object_bucket.performance_monitoring_bucket.endpoint
}

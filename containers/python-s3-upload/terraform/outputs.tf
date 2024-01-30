output "bucket_name" {
  value = scaleway_object_bucket.main.name
}

output "endpoint" {
  value = scaleway_container.main.domain_name
}

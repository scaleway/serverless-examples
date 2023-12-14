
output "endpoint" {
  value = scaleway_container.inference.domain_name
}

output "training_job_id" {
  value = scaleway_job_definition.training.id
}

output "fetch_data_job_id" {
  value = scaleway_job_definition.fetch_data.id
}

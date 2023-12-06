
output "endpoint" {
  value = scaleway_container.inference.domain_name
}

output "training_job_id" {
  value = scaleway_job_definition.training.id
}

output "data_job_id" {
  value = scaleway_job_definition.data.id
}

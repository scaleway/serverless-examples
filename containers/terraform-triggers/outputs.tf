output "public_endpoint" {
  value = scaleway_container.public.domain_name
}

output "public_queue" {
  value = scaleway_mnq_queue.public.sqs[0].url
}

output "private_endpoint" {
  value = scaleway_container.private.domain_name
}

output "private_queue" {
  value = scaleway_mnq_queue.private.sqs[0].url
}

output "sqs_admin_access_key" {
  value = scaleway_mnq_credential.main.sqs_sns_credentials[0].access_key
}

output "sqs_admin_secret_key" {
  value     = scaleway_mnq_credential.main.sqs_sns_credentials[0].secret_key
  sensitive = true
}

output "cockpit_logs_public_container" {
  value = "https://${var.project_id}.dashboard.obs.fr-par.scw.cloud/d/scw-serverless-containers-logs/serverless-containers-logs?orgId=1&var-container_name=${split(".", scaleway_container.public.domain_name)[0]}&var-logs=Scaleway%20Logs"
}

output "cockpit_logs_private_container" {
  value = "https://${var.project_id}.dashboard.obs.fr-par.scw.cloud/d/scw-serverless-containers-logs/serverless-containers-logs?orgId=1&var-container_name=${split(".", scaleway_container.private.domain_name)[0]}&var-logs=Scaleway%20Logs"
}

output "public_endpoint" {
  value = scaleway_container.public.domain_name
}

output "public_queue" {
  value = scaleway_mnq_sqs_queue.public.url
}

output "private_endpoint" {
  value = scaleway_container.private.domain_name
}

output "private_queue" {
  value = scaleway_mnq_sqs_queue.private.url
}

output "sqs_admin_access_key" {
  value = scaleway_mnq_sqs_credentials.main.access_key
  sensitive = true
}

output "sqs_admin_secret_key" {
  value     = scaleway_mnq_sqs_credentials.main.secret_key
  sensitive = true
}

output "public_subject" {
  value = local.public_nats_subject
}

output "private_subject" {
  value = local.private_nats_subject
}

output "nats_creds_file" {
  value = abspath(local_sensitive_file.nats.filename)
}

output "cockpit_logs_public_container" {
  value = "https://${var.project_id}.dashboard.obs.fr-par.scw.cloud/d/scw-serverless-containers-logs/serverless-containers-logs?orgId=1&var-container_name=${split(".", scaleway_container.public.domain_name)[0]}&var-logs=Scaleway%20Logs"
}

output "cockpit_logs_private_container" {
  value = "https://${var.project_id}.dashboard.obs.fr-par.scw.cloud/d/scw-serverless-containers-logs/serverless-containers-logs?orgId=1&var-container_name=${split(".", scaleway_container.private.domain_name)[0]}&var-logs=Scaleway%20Logs"
}

output "public-endpoint" {
  value = scaleway_container.public.domain_name
}

output "public-queue" {
  value = scaleway_mnq_queue.public.sqs[0].url
}

output "private-endpoint" {
  value = scaleway_container.private.domain_name
}

output "private-queue" {
  value = scaleway_mnq_queue.private.sqs[0].url
}

output "sqs_admin_access_key" {
  value = scaleway_mnq_credential.main.sqs_sns_credentials[0].access_key
}

output "sqs_admin_secret_key" {
  value     = scaleway_mnq_credential.main.sqs_sns_credentials[0].secret_key
  sensitive = true
}

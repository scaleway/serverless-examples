output "sqs_access_key" {
  value     = scaleway_mnq_sqs_credentials.main.access_key
  sensitive = true
}

output "sqs_secret_key" {
  value     = scaleway_mnq_sqs_credentials.main.secret_key
  sensitive = true
}

output "subject_name" {
  value = local.subject_name
}

output "creds_file" {
  value = local_sensitive_file.nats.filename
}

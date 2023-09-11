output "sqs_access_key" {
  value = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
}

output "sqs_secret_key" {
  value     = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  sensitive = true
}

output "nats_creds_id" {
  value = scaleway_mnq_nats_credentials.main.id
  sensitive = true
}

output "nats_creds" {
  value = scaleway_mnq_nats_credentials.main.file
  sensitive = true
}

output "nats_endpoint" {
  value = scaleway_mnq_nats_account.main.endpoint
  sensitive = true
}


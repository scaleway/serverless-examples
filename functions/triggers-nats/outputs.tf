output "nats_endpoint" {
  value = scaleway_mnq_nats_account.main.endpoint
  sensitive = true
}


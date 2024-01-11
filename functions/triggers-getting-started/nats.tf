resource "scaleway_mnq_nats_account" "main" {
  name = "nats-account"
}

resource "scaleway_mnq_nats_credentials" "main" {
  account_id = scaleway_mnq_nats_account.main.id
  name       = "triggers-nats"
}

resource "local_sensitive_file" "nats" {
  filename = "triggers-nats.creds"
  content  = scaleway_mnq_nats_credentials.main.file
}

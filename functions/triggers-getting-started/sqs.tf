resource "scaleway_mnq_sqs_credentials" "main" {
  name = "triggers-getting-started"
  permissions {
    can_publish = true
    can_receive = true
    can_manage  = true
  }
}

resource "scaleway_mnq_sqs_queue" "main" {
  for_each = local.functions

  name = "factorial-requests-${each.key}"

  access_key = scaleway_mnq_sqs_credentials.main.access_key
  secret_key = scaleway_mnq_sqs_credentials.main.secret_key
}

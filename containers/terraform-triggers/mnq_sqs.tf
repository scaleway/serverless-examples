# Admin credentials used to create the queues, and to send messages (see tests/send_messages.py)
resource "scaleway_mnq_sqs_credentials" "main" {
  permissions {
    can_publish = true
    can_receive = true
    can_manage  = true
  }
}

locals {
  sqs_admin_credentials_access_key = scaleway_mnq_sqs_credentials.main.access_key
  sqs_admin_credentials_secret_key = scaleway_mnq_sqs_credentials.main.secret_key
}

resource "scaleway_mnq_sqs_queue" "public" {
  name = "sqs-queue-public"

  access_key = local.sqs_admin_credentials_access_key
  secret_key = local.sqs_admin_credentials_secret_key
}

resource "scaleway_mnq_sqs_queue" "private" {
  name = "sqs-queue-private"

  access_key = local.sqs_admin_credentials_access_key
  secret_key = local.sqs_admin_credentials_secret_key
}

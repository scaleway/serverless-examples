resource "scaleway_mnq_namespace" "main" {
  protocol = "sqs_sns"
}

# admin credentials used to create the queues, and to send messages (see tests/send_messages.py)
resource "scaleway_mnq_credential" "main" {
  namespace_id = scaleway_mnq_namespace.main.id

  sqs_sns_credentials {
    permissions {
      can_publish = true
      can_receive = true
      can_manage  = true
    }
  }
}

resource "scaleway_mnq_queue" "public" {
  namespace_id = scaleway_mnq_namespace.main.id
  name         = "sqs-queue-public"

  sqs {
    access_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
    secret_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  }
}

resource "scaleway_mnq_queue" "private" {
  namespace_id = scaleway_mnq_namespace.main.id
  name         = "sqs-queue-private"

  sqs {
    access_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
    secret_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  }
}

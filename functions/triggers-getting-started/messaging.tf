
resource "scaleway_mnq_namespace" "main" {
  name     = "triggers-getting-started"
  protocol = "sqs_sns"
}

resource "scaleway_mnq_credential" "main" {
  name         = "triggers-getting-started"
  namespace_id = scaleway_mnq_namespace.main.id
  sqs_sns_credentials {
    permissions {
      can_publish = true
      can_receive = true
      can_manage  = true
    }
  }
}

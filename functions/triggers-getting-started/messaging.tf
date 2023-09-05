
resource "scaleway_mnq_namespace" "main" {
  name     = "serverless-examples"
  protocol = "sqs_sns"
}

resource "scaleway_mnq_credential" "main" {
  name         = "serverless-examples"
  namespace_id = scaleway_mnq_namespace.main.id
  sqs_sns_credentials {
    permissions {
      can_publish = true
      can_receive = true
      can_manage  = true
    }
  }
}

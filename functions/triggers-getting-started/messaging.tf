resource "scaleway_mnq_sqs_credentials" "main" {
  name         = "triggers-getting-started"
  permissions {
    can_publish = true
    can_receive = true
    can_manage  = true
  }
}

resource "scaleway_container_trigger" "public_sqs" {
  container_id = scaleway_container.public.id
  name         = "public-sqs-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.public.name
  }
}

resource "scaleway_container_trigger" "private_sqs" {
  container_id = scaleway_container.private.id
  name         = "private-sqs-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.private.name
  }
}

locals {
  public_nats_subject = "public-nats-subject"
  private_nats_subject = "private-nats-subject"
}

resource "scaleway_container_trigger" "public_nats" {
  container_id = scaleway_container.public.id
  name         = "public-nats-trigger"
  nats {
    account_id = scaleway_mnq_nats_account.main.id
    subject = local.public_nats_subject
  }
}

resource "scaleway_container_trigger" "private_nats" {
  container_id = scaleway_container.private.id
  name         = "private-nats-trigger"
  nats {
    account_id = scaleway_mnq_nats_account.main.id
    subject = local.private_nats_subject
  }
}

resource "scaleway_container_trigger" "public" {
  container_id = scaleway_container.public.id
  name         = "public-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.public.name
  }
}

resource "scaleway_container_trigger" "private" {
  container_id = scaleway_container.private.id
  name         = "private-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.private.name
  }
}

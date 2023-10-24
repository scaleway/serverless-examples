resource "scaleway_container_trigger" "public" {
  container_id = scaleway_container.public.id
  name         = "public-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.public.name
  }
}

# There is a small bug today when multiple triggers are created at the same time
# We'll wait 10 seconds before creating other triggers
resource "time_sleep" "wait_10_seconds_after_public_trigger_creation" {
  depends_on = [scaleway_container_trigger.public]

  create_duration = "10s"
}

resource "scaleway_container_trigger" "private" {
  container_id = scaleway_container.private.id
  name         = "private-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.private.name
  }

  depends_on = [time_sleep.wait_10_seconds_after_public_trigger_creation]
}

resource "scaleway_container_namespace" "main" {
  name        = "serverless-examples"
  description = "Serverless examples"
}

# This is just to make sure the image is available before creating the containers
# Sometimes, the image is pushed from Terraform perspective but it takes a few seconds for it to be really available
resource "time_sleep" "wait_10_seconds_after_pushing_image" {
  depends_on = [docker_registry_image.main]

  create_duration = "10s"
}

resource "scaleway_container" "public" {
  name           = "example-public-container"
  description    = "Public example container"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port           = 8080
  cpu_limit      = 500
  memory_limit   = 1024
  min_scale      = 0
  max_scale      = 1
  privacy        = "public"
  protocol       = "http1"
  deploy         = true
  depends_on     = [time_sleep.wait_10_seconds_after_pushing_image]
}

resource "scaleway_container" "private" {
  name           = "example-private-container"
  description    = "Private example container"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port           = 8080
  cpu_limit      = 500
  memory_limit   = 1024
  min_scale      = 0
  max_scale      = 1
  privacy        = "private"
  protocol       = "http1"
  deploy         = true
  depends_on     = [time_sleep.wait_10_seconds_after_pushing_image]
}

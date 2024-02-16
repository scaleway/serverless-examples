resource scaleway_container_namespace main {
  name = "mongo-example"
}

resource scaleway_container main {
  name = "mongo-example"
  description = "Example using custom image to run Bash script"
  namespace_id = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port = 8080
  cpu_limit = 1000
  memory_limit = 1024
  min_scale = 0
  max_scale = 1
  privacy = "public"
  deploy = true
  secret_environment_variables = {
    "MONGO_HOSTNAME" = var.mongo_hostname
    "MONGO_USERNAME" = var.mongo_username
    "MONGO_PASSWORD" = var.mongo_password
  }
}

resource scaleway_container_namespace main {
  name = "serverless-example-ns"
  description = "Namespace managed by terraform"
}

resource scaleway_container main {
  name = "serverless-example-container"
  description = "NGINX container deployed with terraform"
  namespace_id = scaleway_container_namespace.main.id
  registry_image = "docker.io/library/nginx:latest"
  port = 80
  cpu_limit = 1120
  memory_limit = 1024
  min_scale = 0
  max_scale = 1
  privacy = "public"
  protocol = "http1"
  deploy = true
}

resource scaleway_container_namespace main {
  name = "python-s3-example"
}

resource scaleway_container main {
  name = "python-s3-example"
  description = "S3 file uploader"
  namespace_id = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port = 80
  cpu_limit = 1000
  memory_limit = 1024
  min_scale = 0
  max_scale = 1
  privacy = "public"
  deploy = true
  environment_variables = {
    "BUCKET_NAME" = scaleway_object_bucket.main.name
  }
  secret_environment_variables = {
    "ACCESS_KEY" = var.access_key
    "SECRET_KEY" = var.secret_key
  }
}

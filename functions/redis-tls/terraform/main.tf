resource "scaleway_function_namespace" "weather_redis_namespace" {
  name        = "weather-redis"
  description = "Part of serverless examples. Connects to Redis with TLS."

  secret_environment_variables = {
    "REDIS_CERT"     = scaleway_redis_cluster.weather_store.certificate
    "REDIS_PASSWORD" = var.redis_password
  }

  environment_variables = {
    "REDIS_URL"  = scaleway_redis_cluster.weather_store.public_network[0].ips[0]
    "REDIS_USER" = var.redis_user
  }
}

data "pypi_requirements_file" "requirements" {
  requirements_file = "${path.module}/../requirements.txt"
  output_dir        = "${path.module}/../package"
}

data "archive_file" "weather_archive" {
  type             = "zip"
  source_dir       = "${path.module}/../"
  excludes         = ["${path.module}/../terraform/**"]
  output_file_mode = "0666"
  output_path      = "${path.module}/files/package.zip"

  depends_on = [
    data.pypi_requirements_file.requirements
  ]
}

resource "scaleway_function" "weather_redis" {
  namespace_id = scaleway_function_namespace.weather_redis_namespace.id
  name         = "weather-to-redis"
  runtime      = "python311"
  handler      = "weather_to_redis.handle"
  privacy      = "public"
  zip_file     = data.archive_file.weather_archive.output_path
  zip_hash     = data.archive_file.weather_archive.output_sha
  deploy       = true
}

resource "scaleway_function_cron" "weather_daily" {
  function_id = scaleway_function.weather_redis.id
  args = ""
  // Everyday at 8:00 UTC
  schedule = "0 8 * * *"
}
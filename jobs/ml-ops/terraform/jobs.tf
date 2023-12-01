resource scaleway_job_definition data {
  name = "data"
  cpu_limit = 1000
  memory_limit = 1024
  image_uri = docker_image.data.name
  timeout = "10m"

  env = {
    "S3_BUCKET_NAME": scaleway_object_bucket.main.name,
    "S3_URL": var.s3_url,
    "SCW_ACCESS_KEY": var.access_key,
    "SCW_SECRET_KEY": var.secret_key,
    "SCW_REGION": var.region
  }
}

resource scaleway_job_definition training {
  name = "training"
  cpu_limit = 6000
  memory_limit = 4096
  image_uri = docker_image.training.name
  timeout = "10m"

  env = {
    "S3_BUCKET_NAME": scaleway_object_bucket.main.name,
    "S3_URL": var.s3_url,
    "SCW_ACCESS_KEY": var.access_key,
    "SCW_SECRET_KEY": var.secret_key,
    "SCW_REGION": var.region,
  }
}

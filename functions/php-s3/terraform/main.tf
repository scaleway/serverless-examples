variable "project_id" {
  description = "Scaleway project ID"
  type        = string
}

variable "access_key" {
  description = "Scaleway access key"
  type        = string
}

variable "secret_key" {
  description = "Scaleway secret key"
  type        = string
  sensitive   = true
}

# Randomness for bucket name uniqueness
resource "random_string" "suffix" {
  length = 8
  special = false
}

# Terraform provider
terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
      version = ">= 2.8.0"
    }
  }

  required_version = ">= 0.13"

  backend "local" {
    path = "state"
  }
}

provider "scaleway" {
  zone   = "fr-par-1"
  region = "fr-par"

  access_key = var.access_key
  secret_key = var.secret_key
  project_id = var.project_id
}

# S3 bucket
resource "scaleway_object_bucket" "main" {
  name = "php-s3-output-${random_string.suffix.result}"
  project_id  = var.project_id
  region      = "fr-par"
}

# Function zip
data "archive_file" "main" {
  type             = "zip"
  source_dir       = "${path.module}/../function"
  excludes         = ["composer.lock", "function.zip", "vendor"]
  output_file_mode = "0666"
  output_path      = "${path.module}/../function/function.zip"
}

# Function namespace and function
resource "scaleway_function_namespace" "main" {
  project_id  = var.project_id
  region      = "fr-par"
  name        = "php-example-namespace"
}

resource "scaleway_function" "main" {
  name           = "php-s3-function"
  description    = "PHP example working with S3"
  namespace_id   = scaleway_function_namespace.main.id
  runtime        = "php82"
  handler        = "handler.run"
  min_scale      = 0
  max_scale      = 2
  zip_file       = data.archive_file.main.output_path
  zip_hash       = data.archive_file.main.output_sha
  privacy        = "public"
  deploy         = true
  memory_limit   = 512

  environment_variables = {
    "S3_ENDPOINT" = "https://s3.fr-par.scw.cloud"
    "S3_REGION" = "fr-par"
    "S3_BUCKET" = scaleway_object_bucket.main.name
  }

  secret_environment_variables = {
    "S3_ACCESS_KEY" = var.access_key
    "S3_SECRET_KEY" = var.secret_key
  }

  timeouts {
    create = "15m"
    update = "15m"
  }
}

output function_url {
  value = scaleway_function.main.domain_name
}

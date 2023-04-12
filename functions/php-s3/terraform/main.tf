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
}

# S3 bucket
resource "scaleway_object_bucket" "example_bucket" {
  name = "php-s3-output"
}

# Function zip
data "archive_file" "func_archive" {
  type             = "zip"
  source_dir       = "${path.module}/../function"
  excludes         = ["composer.lock", "function.zip", "vendor"]
  output_file_mode = "0666"
  output_path      = "${path.module}/../function/function.zip"
}

# Function namespace and function
resource "scaleway_function_namespace" "function_ns" {
  project_id  = var.project_id
  region      = "fr-par"
  name        = "php-example-namespace"
}

resource "scaleway_function" "php_function" {
  name           = "php-s3-function"
  description    = "PHP example working with S3"
  namespace_id   = scaleway_function_namespace.function_ns.id
  runtime        = "php82"
  handler        = "handler.run"
  min_scale      = 0
  max_scale      = 2
  zip_file       = data.archive_file.func_archive.output_path
  zip_hash       = data.archive_file.func_archive.output_sha
  privacy        = "public"
  deploy         = true

  environment_variables = {
    "S3_ENDPOINT" = "https://s3.fr-par.scw.cloud"
    "S3_REGION" = "fr-par"
    "S3_BUCKET" = scaleway_object_bucket.example_bucket.name
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

# Template script to curl
resource local_file curl_script {
  filename = "../curl.sh"
  content = templatefile(
    "curl.tftpl", {
      func_url = "${scaleway_function.php_function.domain_name}",
    }
  )
}

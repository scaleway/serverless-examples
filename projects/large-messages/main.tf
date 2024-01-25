terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
    null = {
      source = "hashicorp/null"
    }
    random = {
      source = "hashicorp/random"
    }
    archive = {
      source = "hashicorp/archive"
    }
  }
  required_version = ">= 0.13"
}

variable "scw_access_key_id" {
  type      = string
  sensitive = true
}

variable "scw_secret_access_key" {
  type      = string
  sensitive = true
}

provider "scaleway" {
  zone = "fr-par-1"
}

// Object Bucket

resource "random_id" "bucket" {
  byte_length = 8
}

resource "scaleway_object_bucket" "large_messages" {
  name = "large-messages-${random_id.bucket.hex}"
}

resource "scaleway_object_bucket_acl" "large_messages" {
  bucket = scaleway_object_bucket.large_messages.id
  acl    = "private"
}

output "bucket_name" {
  value       = scaleway_object_bucket.large_messages.name
  description = "Bucket name to use with the producer script"
}

// MNQ Nats

resource "scaleway_mnq_nats_account" "large_messages" {
  name = "nats-acc-large-messages"
}

resource "scaleway_mnq_nats_credentials" "large_messages" {
  name       = "nats-large-messages-creds"
  account_id = scaleway_mnq_nats_account.large_messages.id
}

resource "local_file" "nats_credential" {
  content         = scaleway_mnq_nats_credentials.large_messages.file
  filename        = "large-messages.creds"
  file_permission = 644
}

output "nats_url" {
  value       = scaleway_mnq_nats_account.large_messages.endpoint
  description = "NATS url to use with the producer script"
}

// Function

resource "null_resource" "install_dependencies" {
  provisioner "local-exec" {
    command = <<-EOT
      cd function
      [ -d "./function/package" ] && rm -rf ./package
      PYTHON_VERSION=3.11
      docker run --rm -v $(pwd):/home/app/function --workdir /home/app/function rg.fr-par.scw.cloud/scwfunctionsruntimes-public/python-dep:$PYTHON_VERSION \
        pip3 install --upgrade -r requirements.txt --no-cache-dir --target ./package
      cd ..
    EOT
  }

  triggers = {
    hash = filesha256("./function/handler/large_messages.py")
  }
}

data "archive_file" "function_zip" {
  type        = "zip"
  source_dir  = "./function"
  output_path = "./function.zip"

  depends_on = [null_resource.install_dependencies]
}

resource "scaleway_function_namespace" "large_messages" {
  name        = "large-messages-function"
  description = "Large messages namespace"
}

resource "scaleway_function" "large_messages" {
  namespace_id = scaleway_function_namespace.large_messages.id
  runtime      = "python311"
  handler      = "handler/large_messages.handle"
  privacy      = "private"
  zip_file     = "function.zip"
  zip_hash     = data.archive_file.function_zip.output_sha256
  deploy       = true
  memory_limit = "2048"
  environment_variables = {
    ENDPOINT_URL  = scaleway_object_bucket.large_messages.api_endpoint
    BUCKET_REGION = scaleway_object_bucket.large_messages.region
    BUCKET_NAME   = scaleway_object_bucket.large_messages.name
  }
  secret_environment_variables = {
    ACCESS_KEY_ID     = var.scw_access_key_id
    SECRET_ACCESS_KEY = var.scw_secret_access_key
  }

  depends_on = [data.archive_file.function_zip]
}

resource "scaleway_function_trigger" "large_messages" {
  function_id = scaleway_function.large_messages.id
  name        = "large-messages-trigger"
  nats {
    account_id = scaleway_mnq_nats_account.large_messages.id
    subject    = "large-messages"
  }
}

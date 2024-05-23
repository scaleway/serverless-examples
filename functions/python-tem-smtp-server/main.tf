terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = ">= 2.40"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.4"
    }
  }
  required_version = ">= 1.0"
}

variable "swc_project_id" {
  type = string
}

variable "scw_secret_key" {
  type = string
  sensitive = true
}

data "archive_file" "function" {
  type        = "zip"
  source_file = "${path.module}/handler.py"
  output_path = "${path.module}/function.zip"
}

resource "scaleway_function_namespace" "main" {
  name        = "serverless-with-tem-example"
  description = "Serverless with TEM example"
}

resource "scaleway_function" "main" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "python-tem-smtp-server"
  runtime      = "python311"
  handler      = "handler.handle"
  privacy      = "public"
  zip_file     = data.archive_file.function.output_path
  zip_hash     = data.archive_file.function.output_sha256
  deploy       = true

  secret_environment_variables = {
    TEM_PROJECT_ID = var.swc_project_id
    SECRET_KEY = var.scw_secret_key
  }
}

output "endpoint" {
  value = scaleway_function.main.domain_name
}

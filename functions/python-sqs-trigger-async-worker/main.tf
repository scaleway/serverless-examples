terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = ">= 2.31"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.4"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3.2"
    }
  }
  required_version = ">= 1.0"
}

# Function resources

locals {
  function_folder_path = "${path.module}/function"
}

resource "null_resource" "pip_install" {
  triggers = {
    requirements = filesha256("${local.function_folder_path}/requirements.txt")
  }

  provisioner "local-exec" {
    command = "pip3 install -r ${local.function_folder_path}/requirements.txt --upgrade --target ${local.function_folder_path}/package"
  }
}

data "archive_file" "function" {
  depends_on = [null_resource.pip_install]

  type        = "zip"
  output_path = "${path.module}/function.zip"

  source_dir = local.function_folder_path
}

resource "scaleway_function_namespace" "main" {
  name        = "serverless-examples"
  description = "Serverless examples"
}

resource "scaleway_function" "front" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "python-sqs-trigger-async-front"
  runtime      = "python311"
  memory_limit = 256
  handler      = "handler.handle_front"
  privacy      = "public"
  zip_file     = data.archive_file.function.output_path
  zip_hash     = data.archive_file.function.output_sha256
  deploy       = true

  min_scale = 0
  max_scale = 1

  secret_environment_variables = {
    SQS_ACCESS_KEY = scaleway_mnq_sqs_credentials.main.access_key
    SQS_SECRET_KEY = scaleway_mnq_sqs_credentials.main.secret_key
    SQS_ENDPOINT   = replace(scaleway_mnq_sqs_queue.main.endpoint, "{region}", scaleway_mnq_sqs_queue.main.region)
    SQS_QUEUE_URL  = scaleway_mnq_sqs_queue.main.url
    SQS_REGION     = scaleway_mnq_sqs_queue.main.region
  }
}

resource "scaleway_function" "worker" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "python-sqs-trigger-async-worker"
  runtime      = "python311"
  memory_limit = 256
  handler      = "handler.handle_worker"
  privacy      = "public"
  zip_file     = data.archive_file.function.output_path
  zip_hash     = data.archive_file.function.output_sha256
  deploy       = true

  min_scale = 0
  max_scale = 1
}

# SQS trigger resources
resource "scaleway_mnq_sqs_credentials" "main" {
  name = "serverless-examples-python-sqs-trigger"

  permissions {
    can_publish = true
    can_receive = true
    can_manage  = true
  }
}

resource "scaleway_mnq_sqs_queue" "main" {
  name = "python-sqs-async-worker"

  access_key = scaleway_mnq_sqs_credentials.main.access_key
  secret_key = scaleway_mnq_sqs_credentials.main.secret_key
}

resource "scaleway_function_trigger" "main" {
  function_id = scaleway_function.worker.id
  name        = "my-trigger"
  sqs {
    queue = scaleway_mnq_sqs_queue.main.name
  }
}

output "front_function_endpoint" {
  value = scaleway_function.front.domain_name
}

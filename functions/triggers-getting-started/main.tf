locals {
  functions = {
    "go" = {
      path    = "go"
      runtime = "go120"
      handler = "Handle"
    }
    "node" = {
      path    = "node"
      runtime = "node20"
      handler = "handler.handle"
    }
    "php" = {
      path    = "php"
      runtime = "php82"
      handler = "handler.handle"
    }
    "python" = {
      path    = "python"
      runtime = "python311"
      handler = "handler.handler"
    }
    "rust" = {
      path    = "rust"
      runtime = "rust165"
      handler = "handler"
    }
  }
}

resource "scaleway_function_namespace" "main" {
  name = "triggers-getting-started"
}

data "archive_file" "function" {
  for_each = local.functions

  type        = "zip"
  source_dir  = "${path.module}/${each.value.path}"
  output_path = "${path.module}/${each.value.path}.zip"
}

resource "scaleway_function" "main" {
  for_each = local.functions

  namespace_id = scaleway_function_namespace.main.id
  name         = each.key
  runtime      = each.value.runtime
  handler      = each.value.handler

  // As of 2023/09/04, only public functions can be triggered by a message queue
  privacy = "public"

  zip_file = data.archive_file.function[each.key].output_path
  zip_hash = data.archive_file.function[each.key].output_sha256
  deploy   = true

  memory_limit = 2048 // 1120 mVCPUs

  min_scale = 0
  max_scale = 1
}

resource "scaleway_mnq_queue" "main" {
  for_each = local.functions

  namespace_id = scaleway_mnq_namespace.main.id
  name         = "factorial-requests-${each.key}"

  sqs {
    access_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
    secret_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  }
}

resource "scaleway_function_trigger" "main" {
  for_each = local.functions

  function_id = scaleway_function.main[each.key].id
  name        = "on-factorial-request"
  sqs {
    namespace_id = scaleway_mnq_namespace.main.id
    queue        = scaleway_mnq_queue.main[each.key].name
  }
}

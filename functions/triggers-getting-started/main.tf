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
  subject_name = "triggers-nats-factorial"
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

  privacy = "public"

  zip_file = data.archive_file.function[each.key].output_path
  zip_hash = data.archive_file.function[each.key].output_sha256
  deploy   = true

  memory_limit = 512 # MB / 280 mVCPU

  min_scale = 0
}

resource "scaleway_function_trigger" "main_sqs" {
  for_each = local.functions

  function_id = scaleway_function.main[each.key].id
  name        = "sqs-on-factorial-request"
  sqs {
    queue = scaleway_mnq_sqs_queue.main[each.key].name
  }
}

resource "scaleway_function_trigger" "main_nats" {
  for_each = local.functions

  function_id = scaleway_function.main[each.key].id
  name        = "nats-on-factorial-request"
  nats {
    account_id = scaleway_mnq_nats_account.main.id
    subject    = local.subject_name
  }
}

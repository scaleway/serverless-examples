resource "scaleway_function_namespace" "main" {
  name = "triggers-nats"
}

data "archive_file" "source_zip" {
  type        = "zip"
  source_dir  = "${path.module}/function"
  output_path = "${path.module}/files/function.zip"
}

resource "scaleway_function" "main" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "nats-target"
  runtime      = "python310"
  handler      = "handler.handle"

  zip_file     = data.archive_file.source_zip.output_path
  zip_hash     = filesha256(data.archive_file.source_zip.output_path)

  privacy = "public"

  deploy   = true

  memory_limit = 512 # MB / 280 mVCPU

  min_scale = 0
}


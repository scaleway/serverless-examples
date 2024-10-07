locals {
  project_name = "node-terraform-rss"
  description  = "A simple RSS feed reader that filters on selected topics."
}

resource "scaleway_function_namespace" "main" {
  name        = local.project_name
  description = local.description
}

resource "null_resource" "npm_install" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "npm install --prefix ${path.module}/function"
  }
}

data "archive_file" "source_zip" {
  type        = "zip"
  source_dir  = "${path.module}/function"
  output_path = "${path.module}/files/function.zip"

  depends_on = [null_resource.npm_install]
}

resource "scaleway_function" "main" {
  namespace_id = scaleway_function_namespace.main.id

  name        = local.project_name
  description = local.description

  runtime = "node22"
  handler = "index.handler"
  privacy = "public"

  zip_file = data.archive_file.source_zip.output_path
  zip_hash = filesha256(data.archive_file.source_zip.output_path)

  environment_variables = {
    "SOURCE_FEED_URL"   = "https://lobste.rs/rss"
    "WORTHWHILE_TOPICS" = "nixos, nix, serverless, terraform"
  }

  deploy = true
}

output "function_url" {
  value = "https://${scaleway_function.main.domain_name}"
}

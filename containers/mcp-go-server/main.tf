###############################################################################
# Terraform configuration for the mcp-go-server Serverless Container example.
#
# Deploys a lightweight Go-based MCP (Model Context Protocol) server on Scaleway
# Serverless Containers. The Docker image is built locally and pushed to the
# Scaleway Container Registry, then referenced by the container.
###############################################################################

terraform {
  required_version = ">= 1.5"

  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "~> 2.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Credentials (SCW_ACCESS_KEY, SCW_SECRET_KEY) and the target project
# (SCW_DEFAULT_PROJECT_ID) are read from the environment by the provider.
provider "scaleway" {
  region = "fr-par"
}

locals {
  name_prefix    = "mcp-go-server"
  tags           = ["serverless-examples", "mcp", "go", "terraform"]
  image_name     = "mcp-server"
  image_tag      = "latest"
  container_port = 8080
  min_scale      = 0
  max_scale      = 5
}

# Scaleway registry namespace names are unique per-region,
# so we add a short random suffix to avoid name collisions.
resource "random_string" "registry_suffix" {
  length  = 4
  lower   = true
  upper   = false
  numeric = true
  special = false
}

resource "scaleway_registry_namespace" "main" {
  name = "${local.name_prefix}-registry-${random_string.registry_suffix.result}"
}

resource "scaleway_container_namespace" "main" {
  name        = "${local.name_prefix}-ns"
  description = "Namespace for the MCP Go server example"
  tags        = local.tags
}

# Build and push the Docker image to the Scaleway Container Registry.
# Note: in a production CI/CD pipeline, you would typically build and push the
# image outside of Terraform. This null_resource keeps the example self-contained.
resource "null_resource" "build_and_push" {
  depends_on = [scaleway_registry_namespace.main]

  triggers = {
    image_tag = timestamp() # Force rebuild on every apply
  }

  provisioner "local-exec" {
    command = <<-EOT
      set -euo pipefail
      # $${...} is escaped so the shell expands SCW_SECRET_KEY from the
      # environment at apply time (Terraform otherwise treats it as a reference).
      docker login "${scaleway_registry_namespace.main.endpoint}" -u nologin -p "$${SCW_SECRET_KEY}"
      docker build --platform linux/amd64 -t "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}" .
      docker push "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}"
    EOT
  }
}

resource "scaleway_container" "main" {
  depends_on = [null_resource.build_and_push]

  name           = local.name_prefix
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}"
  port           = local.container_port
  tags           = local.tags

  cpu_limit          = 1000
  memory_limit_bytes = 1024 * 1024 * 1024 # 1 GiB
  min_scale          = local.min_scale
  max_scale          = local.max_scale

  # WARN: Set privacy = "private" to enforces IAM token authentication via the API Gateway,
  # so only callers with a valid X-Auth-Token header reach the container.
  privacy  = "public"
}

output "container_domain" {
  description = "Public HTTPS endpoint of the deployed MCP server."
  value       = scaleway_container.main.public_endpoint
}

output "mcp_endpoint" {
  description = "Full URL of the /mcp JSON-RPC endpoint."
  value       = "${scaleway_container.main.public_endpoint}/mcp"
}

output "health_endpoint" {
  description = "Health check endpoint."
  value       = "${scaleway_container.main.public_endpoint}/health"
}

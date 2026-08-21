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

provider "scaleway" {
  zone       = "fr-par-1"
  region     = "fr-par"
  access_key = var.scw_access_key
  secret_key = var.scw_secret_key
  project_id = var.project_id
}

variable "scw_access_key" {
  description = "Scaleway IAM access key. Read from SCW_ACCESS_KEY env var if unset."
  type        = string
  default     = null
  sensitive   = true
}

variable "scw_secret_key" {
  description = "Scaleway IAM secret key. Read from SCW_SECRET_KEY env var if unset."
  type        = string
  default     = null
  sensitive   = true
}

variable "project_id" {
  description = "Scaleway project ID. Uses the default project if null."
  type        = string
  default     = null
}

variable "application_name" {
  description = "Value exposed by the get_app_name MCP tool (SCW_APPLICATION_NAME env var)."
  type        = string
  default     = "scaleway-mcp-go-server"
}

variable "container_port" {
  description = "Port the MCP server listens on (matches Dockerfile EXPOSE)."
  type        = number
  default     = 8080
}

variable "min_scale" {
  description = "Minimum number of container instances (0 = scale to zero)."
  type        = number
  default     = 0
}

variable "max_scale" {
  description = "Maximum number of container instances."
  type        = number
  default     = 5
}

locals {
  name_prefix = "mcp-go-server"
  tags        = ["serverless-examples", "mcp", "go", "terraform"]
  image_name  = "mcp-server"
  image_tag   = "latest"
}

# Scaleway registry namespace names are unique per-region,
# so we add a short random suffix to avoid name collisions.
resource "random_string" "registry_suffix" {
  length  = 4
  lower   = true
  special = false
}

resource "scaleway_registry_namespace" "main" {
  name       = "${local.name_prefix}-registry-${random_string.registry_suffix.result}"
  project_id = var.project_id
}

resource "scaleway_container_namespace" "main" {
  name        = "${local.name_prefix}-ns"
  description = "Namespace for the MCP Go server example"
  project_id  = var.project_id
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
      docker login "${scaleway_registry_namespace.main.endpoint}" -u nologin -p "${var.scw_secret_key}"
      docker build -t "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}" .
      docker push "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}"
    EOT
  }
}

resource "scaleway_container" "main" {
  depends_on = [null_resource.build_and_push]

  name           = local.name_prefix
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "${scaleway_registry_namespace.main.endpoint}/${local.image_name}:${local.image_tag}"
  port           = var.container_port
  tags           = local.tags

  cpu_limit          = 1000
  memory_limit_bytes = 1024 * 1024 * 1024 # 1 GiB
  min_scale          = var.min_scale
  max_scale          = var.max_scale

  # privacy = "private" enforces IAM token authentication via the API Gateway,
  # so only callers with a valid X-Auth-Token header reach the container.
  privacy  = "private"
  protocol = "http1"

  environment_variables = {
    PORT                 = tostring(var.container_port)
    SCW_APPLICATION_NAME = var.application_name
  }
}

output "container_domain" {
  description = "Public HTTPS endpoint of the deployed MCP server."
  value       = "https://${scaleway_container.main.public_endpoint}"
}

output "mcp_endpoint" {
  description = "Full URL of the /mcp JSON-RPC endpoint."
  value       = "https://${scaleway_container.main.public_endpoint}/mcp"
}

output "health_endpoint" {
  description = "Health check endpoint."
  value       = "https://${scaleway_container.main.public_endpoint}/health"
}

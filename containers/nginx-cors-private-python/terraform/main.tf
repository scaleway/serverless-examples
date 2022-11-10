variable "project_id" {
  type        = string
  description = ""
}

# Terraform provider
terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
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

# Container namespace and containers
# https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/container
resource "scaleway_container_namespace" "main" {
  project_id  = var.project_id
  region      = "fr-par"
  name        = "cors-demo-namespace"
}

resource "scaleway_container" "server_container" {
  name           = "server-container"
  description    = "Container for the server"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "rg.fr-par.scw.cloud/cors-demo/server:0.0.1"
  port           = 8080
  min_scale      = 2
  max_scale      = 2
  privacy        = "private"
  deploy         = true

  timeouts {
    create = "15m"
    update = "15m"
  }
}

resource "scaleway_container" "gateway_container" {
  name           = "gateway-container"
  description    = "Container for the gateway"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "rg.fr-par.scw.cloud/cors-demo/gateway:0.0.1"
  port           = 8080
  min_scale      = 1
  max_scale      = 1
  privacy        = "public"
  deploy         = true

  environment_variables = {
    "SERVER_CONTAINER_URL" = scaleway_container.server_container.domain_name
  }

  timeouts {
    create = "15m"
    update = "15m"
  }
}

# Token for private container
# https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/container_token
resource scaleway_container_token server {
  container_id = scaleway_container.server_container.id
}

resource local_sensitive_file token {
  filename = "../token"
  content = scaleway_container_token.server.token
}

# Template HTML file
data template_file html_data {
  template = "${file("../index.tftpl")}"
  vars = {
    gateway_func_url = "${scaleway_container.gateway_container.domain_name}",
    server_auth_token = "${scaleway_container_token.server.token}",
  }
}

resource local_sensitive_file html {
  filename = "../index.html"
  content = templatefile(
    "../index.tftpl", {
      gateway_func_url = "${scaleway_container.gateway_container.domain_name}",
      server_auth_token = "${scaleway_container_token.server.token}",
    }
  )
}

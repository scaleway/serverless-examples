locals {
  name_prefix = "mongodb-example"
  tags        = ["serverless-examples", "mongodb", "terraform"]
}

resource "scaleway_account_project" "main" {
  name = "${local.name_prefix}-project"
}

resource "scaleway_vpc" "main" {
  project_id = scaleway_account_project.main.id

  name = "${local.name_prefix}-vpc"
  tags = local.tags
}

resource "scaleway_vpc_private_network" "main" {
  project_id = scaleway_account_project.main.id
  vpc_id = scaleway_vpc.main.id

  name = "${local.name_prefix}-private-network"
  tags = local.tags
}

resource "scaleway_mongodb_instance" "main" {
  project_id = scaleway_account_project.main.id

  name = "${local.name_prefix}-instance"
  tags = local.tags

  version = "7.0.12"

  node_type         = "MGDB-PLAY2-NANO"
  node_number       = 1
  volume_size_in_gb = 10

  public_network {}

  private_network {
    pn_id = scaleway_vpc_private_network.main.id
  }

  user_name = var.mongodb_admin_username
  password  = var.mongodb_admin_password

  is_snapshot_schedule_enabled = false
}

resource "scaleway_mongodb_user" "app_user" {
  instance_id = scaleway_mongodb_instance.main.id
  name        = var.mongodb_username
  password    = var.mongodb_password

  roles {
    role          = "read_write"
    any_database =  true
  }
}

# Scaleway registry namespace names are unique per-region, 
# so we add a random suffix to avoid name collisions.
resource "random_string" "unique_namespace_suffix" {
  length = 4
  lower  = true
  special = false
}

resource "scaleway_registry_namespace" "main" {
  project_id = scaleway_account_project.main.id

  name = "${local.name_prefix}-registry-${random_string.unique_namespace_suffix.result}"
}

resource "scaleway_container_namespace" "main" {
  project_id = scaleway_account_project.main.id

  name        = "${local.name_prefix}-namespace"
  description = "Namespace for MongoDB client"
  tags        = local.tags
}

locals {
  mongodb_pn           = scaleway_mongodb_instance.main.private_network[0]
  mongodb_private_endpoint   = local.mongodb_pn.dns_records[0]
  mongodb_private_port = local.mongodb_pn.port
}

// Note: this is somewhat hacky, in a more real world scenario you'd likely want to
// build and push the Docker image outside of Terraform, e.g. in a CI/CD pipeline.
resource "null_resource" "build_and_push_docker_image" {
  depends_on = [scaleway_registry_namespace.main]
  
  provisioner "local-exec" {
    command = <<EOT
      docker build -t ${scaleway_registry_namespace.main.endpoint}/mongodb-example:latest app
      docker push ${scaleway_registry_namespace.main.endpoint}/mongodb-example:latest
    EOT
  }

  triggers = {
    image_tag = timestamp() // Force rebuild on every apply
  }
}

resource "scaleway_container" "main" {
  depends_on = [null_resource.build_and_push_docker_image]
  namespace_id = scaleway_container_namespace.main.id

  name = "${local.name_prefix}-container"
  tags = local.tags

  registry_image = "${scaleway_registry_namespace.main.endpoint}/mongodb-example:latest"
  port           = 8080

  cpu_limit    = 1000
  memory_limit = 1024

  private_network_id = scaleway_vpc_private_network.main.id

  min_scale = 1
  max_scale = 1

  environment_variables = {
    "MONGODB_CERT" : scaleway_mongodb_instance.main.tls_certificate,
  }

  secret_environment_variables = {
    "MONGODB_URI" : "mongodb+srv://${var.mongodb_username}:${var.mongodb_password}@${local.mongodb_private_endpoint}"
  }

  deploy = true
}

variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "scw_registry" {
  type = string
  default = "rg.fr-par.scw.cloud"
}

provider "scaleway" {
  zone   = "fr-par-1"
  region = "fr-par"
  access_key = var.access_key
  secret_key = var.secret_key
  project_id = var.project_id
}

provider "docker" {
  registry_auth {
    address  = var.scw_registry
    username = var.access_key
    password = var.secret_key
  }
}

# ----- DOCKER BUILD ------

resource "docker_image" "main" {
  name          = "${scaleway_container_namespace.main.registry_endpoint}/server:0.0.1"
  build {
    context = "docker"
  }
  triggers = {
    dir_sha1 = sha1(join("", [for f in fileset(path.module, "docker/*") : filesha1(f)]))
  }
}

resource "docker_registry_image" "main" {
  name          = docker_image.main.name
  keep_remotely = true
}

# ----- MNQ ------

resource "scaleway_mnq_namespace" "main" {
  protocol = "sqs_sns"
}

resource "scaleway_mnq_credential" "main" {
  namespace_id = scaleway_mnq_namespace.main.id

  sqs_sns_credentials {
    permissions {
      can_publish = true
      can_receive = true
      can_manage  = true
    }
  }
}

resource "scaleway_mnq_queue" "public" {
  namespace_id = scaleway_mnq_namespace.main.id
  name         = "sqs-queue-public"

  sqs {
    access_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
    secret_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  }
}

resource "scaleway_mnq_queue" "private" {
  namespace_id = scaleway_mnq_namespace.main.id
  name         = "sqs-queue-private"

  sqs {
    access_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.access_key
    secret_key = scaleway_mnq_credential.main.sqs_sns_credentials.0.secret_key
  }
}

# ----- CONTAINERS ------

resource "scaleway_container_namespace" "main" {
  name = "serverless-examples"
  description = "Serverless examples"
}

# This is just to make sure the image is available before creating the containers
resource "time_sleep" "wait_10_seconds_after_pushing_image" {
  depends_on = [docker_registry_image.main]

  create_duration = "10s"
}

resource "scaleway_container" "public" {
  name = "example-public-container"
  description = "Public example container"
  namespace_id = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port = 8080
  cpu_limit = 500
  memory_limit = 1024
  min_scale = 0
  max_scale = 1
  privacy = "public"
  protocol = "http1"
  deploy = true
  depends_on = [time_sleep.wait_10_seconds_after_pushing_image]
}

resource "scaleway_container" "private" {
  name = "example-private-container"
  description = "Private example container"
  namespace_id = scaleway_container_namespace.main.id
  registry_image = docker_image.main.name
  port = 8080
  cpu_limit = 500
  memory_limit = 1024
  min_scale = 0
  max_scale = 1
  privacy = "private"
  protocol = "http1"
  deploy = true
  depends_on = [time_sleep.wait_10_seconds_after_pushing_image]
}

# ----- TRIGGERS -----

resource "scaleway_container_trigger" "public" {
  container_id = scaleway_container.public.id
  name = "public-trigger"
  sqs {
    namespace_id = scaleway_mnq_namespace.main.id
    queue = scaleway_mnq_queue.public.name
  }
}

resource "scaleway_container_trigger" "private" {
  container_id = scaleway_container.private.id
  name = "private-trigger"
  sqs {
    namespace_id = scaleway_mnq_namespace.main.id
    queue = scaleway_mnq_queue.private.name
  }
}

# ----- OUTPUTS -----

output "public-endpoint" {
  value = scaleway_container.public.domain_name
}

output "public-queue" {
  value = scaleway_mnq_queue.public.sqs
  sensitive = true
}

output "private-endpoint" {
  value = scaleway_container.private.domain_name
}

output "private-queue" {
  value = scaleway_mnq_queue.private.sqs
  sensitive = true
}

# Configuring Producion environment
resource "scaleway_instance_ip" "public_ip-prod" {
  project_id = var.project_id
}

resource "scaleway_instance_server" "scw-instance-prod" {
  name       = "prod"
  project_id = var.project_id
  type       = "GP1-XS"
  image      = "ubuntu_noble"

  tags = ["terraform instance", "scw-instance", "production"]

  ip_id = scaleway_instance_ip.public_ip-prod.id

  root_volume {
    # The local storage of a DEV1-L instance is 80 GB, subtract 30 GB from the additional l_ssd volume, then the root volume needs to be 50 GB.
    size_in_gb = 200
  }
}


# Configuring Development environment that will be automatically turn off on week-ends and turn on monday mornings
resource "scaleway_instance_ip" "public_ip-dev" {
  project_id = var.project_id
}


resource "scaleway_instance_server" "scw-instance-dev" {
  name       = "dev"
  project_id = var.project_id
  type       = "DEV1-S"
  image      = "ubuntu_noble"

  tags = ["terraform instance", "scw-instance", "dev"]

  ip_id = scaleway_instance_ip.public_ip-dev.id


  root_volume {
    size_in_gb = 80
  }
}

# Creating function code archive that will then be updated
data "archive_file" "source_zip" {
  type        = "zip"
  source_dir  = "${path.module}/function"
  output_path = "${path.module}/files/function.zip"
}

# Create an IAM application to provide the a scoped token for the function
resource "scaleway_iam_application" "instance_power_toggler" {
  name        = "Instance power toggler"
  description = "IAM application for the function to toggle instances"
}

resource "scaleway_iam_policy" "can_manage_instances" {
  name           = "can-manage-instances"
  description    = "policy to manage instances in a specific project"
  application_id = scaleway_iam_application.instance_power_toggler.id
  rule {
    project_ids = [var.project_id]
    # Currently the least possible privilege to manage instances
    permission_set_names = ["InstancesFullAccess"]
  }
}

resource "scaleway_iam_api_key" "instance_toggler_key" {
  application_id = scaleway_iam_application.instance_power_toggler.id
  description    = "API key for the function to toggle instances"
}

# Creating the function
resource "scaleway_function_namespace" "main" {
  name        = "instance-management"
  description = "namespace to gather all functions dedicated to instance management"
  project_id  = var.project_id
}

resource "scaleway_function" "main" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "instancewake"
  runtime      = "python312"
  handler      = "handler.handle"
  privacy      = "public"
  zip_file     = data.archive_file.source_zip.output_path
  zip_hash     = filesha256(data.archive_file.source_zip.output_path)
  deploy       = true
  max_scale    = "5"
  environment_variables = {
    "X-AUTH-TOKEN" = scaleway_iam_api_key.instance_toggler_key.secret_key
  }
}

# Adding a first cron to turn off the instance every friday evening (11:30 pm)
resource "scaleway_function_cron" "turn-off" {
  function_id = scaleway_function.main.id
  schedule    = "30 23 * * 5"
  args = jsonencode({
    "zone" : scaleway_instance_server.scw-instance-dev.zone,
    "server_id" : regex("/([^/]+$)", scaleway_instance_server.scw-instance-dev.id)[0],
    "action" : "poweroff"
    }
  )
}

# Adding a second cron to turn on the instance every monday morning (7:00 am)
resource "scaleway_function_cron" "turn-on" {
  function_id = scaleway_function.main.id
  schedule    = "0 7 * * 1"
  args = jsonencode({
    "zone" : scaleway_instance_server.scw-instance-dev.zone,
    "server_id" : regex("/([^/]+$)", scaleway_instance_server.scw-instance-dev.id)[0],
    "action" : "poweron"
    }
  )
}

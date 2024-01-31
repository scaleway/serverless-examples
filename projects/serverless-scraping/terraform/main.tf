terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

provider "scaleway" {
}

resource "scaleway_account_project" "mnq_tutorial" {
  name = "mnq-tutorial"
}

# ============= SQS ===============
resource "scaleway_mnq_sqs" "main" {
  project_id = scaleway_account_project.mnq_tutorial.id
}

resource "scaleway_mnq_sqs_credentials" "producer_creds" {
  project_id = scaleway_mnq_sqs.main.project_id
  name = "sqs-credentials-producer"

  permissions {
    can_manage  = true
    can_receive = false
    can_publish = true
  }
}

resource "scaleway_mnq_sqs_credentials" "consumer_creds" {
  project_id = scaleway_mnq_sqs.main.project_id
  name = "sqs-credentials-consumer"

  permissions {
    can_manage  = false
    can_receive = true
    can_publish = false
  }
}

resource "scaleway_mnq_sqs_queue" "main" {
  # possible to pass scaleway_mnq_sqs resource directly? To see with devtools
  project_id = scaleway_account_project.mnq_tutorial.id
  name       = "hn-queue"
  access_key = scaleway_mnq_sqs_credentials.producer_creds.access_key
  secret_key = scaleway_mnq_sqs_credentials.producer_creds.secret_key
}

# ============= RDB ===============

resource "random_password" "dev_mnq_pg_exporter_password" {
  length           = 16
  special          = true
  min_numeric      = 1
  min_upper        = 1
  min_lower        = 1
  min_special      = 1
  override_special = "_-"
}

resource "scaleway_rdb_instance" "main" {
  name = "test-rdb"
  project_id   = scaleway_account_project.mnq_tutorial.id
  node_type = "db-dev-s"
  engine = "PostgreSQL-15"
  is_ha_cluster = false
  disable_backup = true 
  user_name = "mnq_initial_user"
  password = random_password.dev_mnq_pg_exporter_password.result
}

output db_password {
  value = random_password.dev_mnq_pg_exporter_password.result
  sensitive = true
}

output db_ip {
  value = scaleway_rdb_instance.main.endpoint_ip
  sensitive = false
}

output db_port {
  value = scaleway_rdb_instance.main.endpoint_port
  sensitive = false
}

resource "scaleway_rdb_database" "main" {
  instance_id = scaleway_rdb_instance.main.id 
  name = "hn-database"
}

resource "scaleway_rdb_user" "worker" {
  instance_id = scaleway_rdb_instance.main.id
  name = "worker"
  password = random_password.dev_mnq_pg_exporter_password.result
  is_admin = false
}

resource "scaleway_rdb_privilege" "mnq_user_role" {
  instance_id = scaleway_rdb_instance.main.id 
  user_name = scaleway_rdb_user.worker.name
  database_name = scaleway_rdb_database.main.name
  permission = "all"
}

# ============= Functions ===============

locals {
  scraper_folder_path = "../scraper"
  consumer_folder_path = "../consumer"
  archives_folder_path = "../archives"
}

resource "scaleway_function_namespace" "mnq_tutorial_namespace" {
  project_id = scaleway_account_project.mnq_tutorial.id
  name        = "mnq-tutorial-namespace"
  description = "Main function namespace"
}

resource "null_resource" "pip_install_scraper" {
  triggers = {
    requirements = filesha256("${local.scraper_folder_path}/requirements.txt")
  }

  provisioner "local-exec" {
    command = "pip3 install -r ${local.scraper_folder_path}/requirements.txt --upgrade --target ${local.scraper_folder_path}/package"
  }
}

data "archive_file" "scraper_archive" {
  depends_on = [ null_resource.pip_install_scraper ]


  type = "zip"
  output_path = "${local.archives_folder_path}/scraper.zip"

  source_dir = local.scraper_folder_path
}

resource "scaleway_function" "scraper" {
  namespace_id = scaleway_function_namespace.mnq_tutorial_namespace.id
  project_id   = scaleway_account_project.mnq_tutorial.id
  name         = "mnq-hn-scraper"
  runtime      = "python311"
  handler      = "handlers/scrape_hn.handle"
  privacy      = "private"
  timeout      = 10
  zip_file     = data.archive_file.scraper_archive.output_path
  zip_hash     = data.archive_file.scraper_archive.output_sha256
  deploy       = true
  environment_variables = {
    QUEUE_URL = scaleway_mnq_sqs_queue.main.url
    SQS_ACCESS_KEY = scaleway_mnq_sqs_credentials.producer_creds.access_key
  }
  secret_environment_variables = {
    SQS_SECRET_ACCESS_KEY = scaleway_mnq_sqs_credentials.producer_creds.secret_key
  }
}

resource "null_resource" "pip_install_consumer" {
  triggers = {
    requirements = filesha256("${local.consumer_folder_path}/requirements.txt")
  }

  provisioner "local-exec" {
    command = "pip3 install -r ${local.consumer_folder_path}/requirements.txt --upgrade --target ${local.consumer_folder_path}/package"
  }
}

data "archive_file" "consumer_archive" {
  depends_on = [ null_resource.pip_install_consumer ]


  type = "zip"
  output_path = "${local.archives_folder_path}/consumer.zip"

  source_dir = local.consumer_folder_path
}

resource "scaleway_function" "consumer" {
  namespace_id = scaleway_function_namespace.mnq_tutorial_namespace.id
  project_id   = scaleway_account_project.mnq_tutorial.id
  name         = "mnq-hn-consumer"
  runtime      = "python311"
  handler      = "handlers/consumer.handle"
  privacy      = "private"
  timeout      = 10
  zip_file     = data.archive_file.consumer_archive.output_path
  zip_hash     = data.archive_file.consumer_archive.output_sha256
  deploy       = true
  max_scale    = 3
  environment_variables = {
    DB_NAME = scaleway_rdb_database.main.name
    DB_HOST = scaleway_rdb_instance.main.load_balancer[0].ip
    DB_PORT = scaleway_rdb_instance.main.load_balancer[0].port
    DB_USER = scaleway_rdb_user.worker.name
  }
  secret_environment_variables = {
    DB_PASSWORD = scaleway_rdb_user.worker.password
  }
}

# ============= Triggers ===============

resource "scaleway_function_cron" "scraper_cron" {
  function_id = scaleway_function.scraper.id 
  schedule = "0,15,30,45 * * * *"
  args = jsonencode({})
}

resource "scaleway_function_trigger" "consumer_sqs_trigger" {
  function_id = scaleway_function.consumer.id 
  name = "hn-sqs-trigger"
  sqs {
    project_id = scaleway_mnq_sqs.main.project_id
    queue = scaleway_mnq_sqs_queue.main.name
  }
}
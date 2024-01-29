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

output db_password {
  value = random_password.dev_mnq_pg_exporter_password.result
  sensitive = true
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

resource "scaleway_rdb_database" "hn-database" {
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
  database_name = scaleway_rdb_database.hn-database.name
  permission = "all"
}

# ============= Functions ===============

resource "scaleway_function_namespace" "mnq_tutorial_namespace" {
  project_id = scaleway_account_project.mnq_tutorial.id
  name        = "mnq-tutorial-namespace"
  description = "Main function namespace"
}

resource "scaleway_function" "scraper" {
  namespace_id = scaleway_function_namespace.mnq_tutorial_namespace.id
  project_id   = scaleway_account_project.mnq_tutorial.id
  name         = "mnq-hn-scraper"
  runtime      = "python311"
  handler      = "handlers/scrape_hn.handle"
  privacy      = "private"
  timeout      = 10
  zip_file     = "../scraper/functions.zip"
  zip_hash     = filesha256("../scraper/functions.zip")
  deploy       = true
  environment_variables = {
    QUEUE_URL = scaleway_mnq_sqs_queue.main.url
    SQS_ACCESS_KEY = scaleway_mnq_sqs_credentials.producer_creds.access_key
  }
  secret_environment_variables = {
    SQS_SECRET_ACCESS_KEY = scaleway_mnq_sqs_credentials.producer_creds.secret_key
  }
}

resource "scaleway_function" "consumer" {
  namespace_id = scaleway_function_namespace.mnq_tutorial_namespace.id
  project_id   = scaleway_account_project.mnq_tutorial.id
  name         = "mnq-hn-consumer"
  runtime      = "python311"
  handler      = "handlers/consumer.handle"
  privacy      = "private"
  timeout      = 10
  zip_file     = "../consumer/functions.zip"
  zip_hash     = filesha256("../consumer/functions.zip")
  deploy       = true
  max_scale    = 3
  environment_variables = {
    DB_NAME = scaleway_rdb_database.hn-database.name
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
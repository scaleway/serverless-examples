terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

# This Terraform configuration deploys a Memos instance on Scaleway.
# Memos is an "open-source, lightweight note-taking solution. The pain-less way to create your meaningful notes. Your Notes, Your Way."
#
# This configuration creates the necessary resources, including:
# - A Scaleway project
# - An IAM application and policy for Serverless SQL Database access
# - A Serverless SQL database (postgres)
# - A Serverless Container to run Memos

# Create a new project "memos"
resource "scaleway_account_project" "project" {
  name = "memos"
}

# Create a new IAM Application called "memos"
resource "scaleway_iam_application" "app" {
  name = "memos"
}

# Create a new IAM Policy that will give the API Key access to the Database
resource "scaleway_iam_policy" "db_access" {
  name           = "memos_policy"
  description    = "access to serverless database in project"
  application_id = scaleway_iam_application.app.id
  rule {
    project_ids          = [scaleway_account_project.project.id]
    permission_set_names = ["ServerlessSQLDatabaseReadWrite"]
  }
}

# Create API Key
resource "scaleway_iam_api_key" "api_key" {
  application_id = scaleway_iam_application.app.id
}

# Create a new SQL Serverless Database in the "memos" project
resource "scaleway_sdb_sql_database" "database" {
  name       = "memos"
  min_cpu    = 0
  max_cpu    = 1
  project_id = scaleway_account_project.project.id
}

locals {
  # This connection string provides access to the Memos database.
  # It will be used to inject in the Serverless Container
  database_connection_string = format("postgres://%s:%s@%s",
    scaleway_iam_application.app.id,
    scaleway_iam_api_key.api_key.secret_key,
    trimprefix(scaleway_sdb_sql_database.database.endpoint, "postgres://"),
  )
}

# Create namespace for the Serverless Container
resource "scaleway_container_namespace" "main" {
  name       = "memos-container"
  project_id = scaleway_account_project.project.id
}

# Create the container with memos image. You can adjust cpu and memory
# depending the traffic on your application.
# Do not change the port
resource "scaleway_container" "main" {
  name           = "memos"
  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "neosmemo/memos:stable"
  port           = 5230
  cpu_limit      = 1000
  memory_limit   = 2048
  min_scale      = 0
  max_scale      = 5
  privacy        = "public"
  protocol       = "http1"
  deploy         = true

  environment_variables = {
    # Specifies the database driver to use.
    "MEMOS_DRIVER" = "postgres"
  }
  secret_environment_variables = {
    # Provides the database connection string for Memos to connect to the database.
    "MEMOS_DSN" = local.database_connection_string
  }
}

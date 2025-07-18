locals {
  name                = "metabase-example"
  db_postgres_version = "15"
  base_tags           = ["metabase", "vpc"]
}

resource "scaleway_vpc" "main" {
  name = local.name
}

resource "scaleway_vpc_private_network" "main" {
  name   = local.name
  vpc_id = scaleway_vpc.main.id
}

resource "scaleway_rdb_instance" "main" {
  name = "db-${local.name}"
  tags = concat(local.base_tags, ["pg-${local.db_postgres_version}"])

  node_type = "db-play2-nano"

  is_ha_cluster = false
  private_network {
    pn_id       = scaleway_vpc_private_network.main.id
    enable_ipam = true
  }

  encryption_at_rest = true
  volume_size_in_gb  = 10
  volume_type        = "sbs_5k"

  engine = "PostgreSQL-${local.db_postgres_version}"

  user_name = var.db_admin_username
  password  = var.db_admin_password
}

resource "scaleway_rdb_database" "main" {
  instance_id = scaleway_rdb_instance.main.id
  name        = local.name
}

resource "scaleway_rdb_user" "main" {
  instance_id = scaleway_rdb_instance.main.id

  name     = var.db_username
  password = var.db_password
}

resource "scaleway_rdb_privilege" "main" {
  instance_id   = scaleway_rdb_instance.main.id
  user_name     = scaleway_rdb_user.main.name
  database_name = scaleway_rdb_database.main.name
  permission    = "all"
}

resource "scaleway_container_namespace" "main" {
  name                     = local.name
  description              = "Namespace for the Metabase container"
  tags                     = local.base_tags
  
  activate_vpc_integration = true
}

locals {
  db_endpoint     = scaleway_rdb_instance.main.private_network[0]
  rdb_instance_id = split("/", scaleway_rdb_instance.main.id)[1] # To remove the `<region>/` prefix
}

resource "scaleway_container" "main" {
  name        = local.name
  description = "Metabase container running in VPC"
  tags        = local.base_tags

  namespace_id   = scaleway_container_namespace.main.id
  registry_image = "metabase/metabase:v0.55.x"

  private_network_id = scaleway_vpc_private_network.main.id

  cpu_limit    = 4000
  memory_limit = 4000
  sandbox      = "v1"

  http_option = "redirected" # Only allow HTTPs traffic
  port        = 3000

  max_scale = 1 # No real need to have more than one instance running
  deploy    = true

  environment_variables = {
    MB_ANON_TRACKING_ENABLED : "false"
    MB_CHECK_FOR_UPDATES : "false"

    MB_JETTY_HOST : "0.0.0.0"

    MB_DB_TYPE : "postgres"
    MB_DB_CONNECTION_TIMEOUT_MS : "2000" # Down from 10s for faster feedback loop

    # Within a private network, we can refer to resources using their internal hostname
    # The format is `<resource_id>.<private_network_name>.internal` or `<resource_name>.<private_network_name>.internal`
    MB_DB_HOST : "${local.rdb_instance_id}.${scaleway_vpc_private_network.main.name}.internal"
    MB_DB_PORT : local.db_endpoint.port

    MB_DB_DBNAME : scaleway_rdb_database.main.name
    MB_DB_USER : scaleway_rdb_user.main.name # Referencing the user directly to create a Terraform dependency
  }

  secret_environment_variables = {
    MB_DB_PASS : var.db_password
  }
}

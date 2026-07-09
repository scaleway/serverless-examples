locals {
  name                = "pg-connection-retry"
  db_postgres_version = "16"
  base_tags           = ["pg-connection-retry", "vpc"]
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
  name        = local.name
  description = "Namespace for the pg-connection-retry container"
  tags        = local.base_tags
}

locals {
  db_endpoint     = scaleway_rdb_instance.main.private_network[0]
  rdb_instance_id = split("/", scaleway_rdb_instance.main.id)[1] # To remove the `<region>/` prefix
}

resource "scaleway_container" "main" {
  name        = local.name
  description = "Node.js container that retries its VPC connection to a Postgres DB"
  tags        = local.base_tags

  namespace_id = scaleway_container_namespace.main.id
  image        = var.registry_image

  private_network_id = scaleway_vpc_private_network.main.id

  cpu_limit          = 1000
  memory_limit_bytes = 1024 * 1024 * 1024 # 1 GB
  sandbox            = "v1"

  https_connections_only = true
  port                   = 8080

  max_scale = 1 # No real need to have more than one instance running

  environment_variables = {
    # Within a private network, we can refer to resources using their internal hostname.
    # The format is `<resource_id>.<private_network_name>.internal`.
    PGHOST = "${local.rdb_instance_id}.${scaleway_vpc_private_network.main.name}.internal"
    PGPORT = local.db_endpoint.port

    PGDATABASE = scaleway_rdb_database.main.name
    PGUSER     = scaleway_rdb_user.main.name # Referencing the user directly to create a Terraform dependency
  }

  secret_environment_variables = {
    PGPASSWORD = var.db_password
  }
}

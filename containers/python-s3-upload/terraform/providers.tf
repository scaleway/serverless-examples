provider "scaleway" {
  zone   = "fr-par-1"
  region = "fr-par"
  access_key = var.access_key
  secret_key = var.secret_key
  project_id = var.project_id
}

provider "docker" {
  host = "unix:///var/run/docker.sock"

  registry_auth {
    address  = scaleway_registry_namespace.main.endpoint
    username = "nologin"
    password = var.secret_key
  }
}

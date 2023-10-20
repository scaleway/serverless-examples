provider "scaleway" {
  zone       = "fr-par-1"
  region     = "fr-par"
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

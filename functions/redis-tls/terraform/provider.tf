terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
    pypi = {
      source  = "jeffwecan/pypi"
      version = "0.0.11"
    }
  }
  required_version = ">= 0.13"
}

provider "scaleway" {
  region = "fr-par"
}

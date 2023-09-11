terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = ">= 2.12"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.4"
    }
  }
  required_version = ">= 0.13"
}

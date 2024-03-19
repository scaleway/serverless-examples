terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
      version = ">= 2.38.2"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
  required_version = ">= 0.13"
}

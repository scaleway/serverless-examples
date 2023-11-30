terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
  required_version = ">= 0.13"
}

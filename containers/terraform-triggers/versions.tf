terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
      version = ">= 0.13"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">=3.0.2"
    }
  }
}

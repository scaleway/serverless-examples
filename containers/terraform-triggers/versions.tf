terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "2.30.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.9.1"
    }
  }
}

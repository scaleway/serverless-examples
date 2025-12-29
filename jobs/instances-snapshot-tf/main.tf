terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "zone" {
  type = string
}

variable "block_volume_id" {
  type = string
}

resource "scaleway_job_definition" "main" {
  name         = "job-example-using-cli"
  cpu_limit    = 140
  memory_limit = 256
  image_uri    = "scaleway/cli:latest"
  command      = "/scw block snapshot create volume-id=${var.block_volume_id} --debug"
  timeout      = "2m"
  project_id   = scaleway_project.snapshot_instance.id

  env = {
    "SCW_DEFAULT_ORGANIZATION_ID" : var.project_id,
    "SCW_DEFAULT_PROJECT_ID" : scaleway_project.snapshot_instance.id,
    "SCW_DEFAULT_REGION" : var.region,
    "SCW_DEFAULT_ZONE" : var.zone
  }

  cron {
    schedule = "0 * * * *" # run at midnight every day
    timezone = "Europe/Paris"
  }
}

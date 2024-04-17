variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "image_version" {
  type = string
  default = "0.0.3"
}

variable "region" {
  type = string
  default = "fr-par"
}

variable "s3_url" {
  type = string
  default = "https://s3.fr-par.scw.cloud"
}

variable "data_fetch_cron_schedule" {
  type = string
}

variable "training_cron_schedule" {
  type = string
  default = "*/15 * * * *"
}

variable "inference_cron_schedule" {
  type = string
  default = "*/25 * * * *"
}

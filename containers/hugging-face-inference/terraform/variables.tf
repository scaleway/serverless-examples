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

variable "inference_cron_schedule" {
  type = string
  default = "*/15 * * * *"
}

variable "hf_model_file_name" {
  type = string
}

variable "hf_model_download_source" {
  type = string
}

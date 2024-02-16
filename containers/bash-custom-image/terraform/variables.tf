variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "mongo_username" {
  type = string
}

variable "mongo_password" {
  type = string
}

variable "mongo_hostname" {
  type = string
}

variable "image_version" {
  type = string
  default = "0.0.1"
}

resource "random_string" "suffix" {
  length  = 8
  upper   = false
  special = false
}

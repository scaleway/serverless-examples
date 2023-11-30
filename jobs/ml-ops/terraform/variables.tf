variable "region" {
  type = string
}

variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "data_file" {
  type        = string
  description = "name data file in data store"
}

variable "model_object" {
  type        = string
  description = "name of model object stored in model registry"
}

variable "image_version" {
  type = string
}


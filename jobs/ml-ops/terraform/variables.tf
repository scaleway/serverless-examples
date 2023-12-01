variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "region" {
  type = string
  default = "fr-par"
}

variable "s3_url" {
  type = string
  default = "https://s3.fr-par.scw.cloud"
}

variable "data_file" {
  type        = string
  description = "name data file in data store"
  default = "bank_telemarketing.csv"
}

variable "model_object" {
  type        = string
  description = "name of model object stored in model registry"
  default = "classifier.pkl"
}

variable "image_version" {
  type = string
  default = "0.0.1"
}


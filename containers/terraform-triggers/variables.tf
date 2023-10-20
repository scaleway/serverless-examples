variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "scw_registry" {
  type    = string
  default = "rg.fr-par.scw.cloud"
}

variable "container_language" {
  type    = string
  default = "python"
}

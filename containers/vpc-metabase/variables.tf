variable "db_admin_username" {
  type    = string
  default = "admin"
}

variable "db_admin_password" {
  type      = string
  sensitive = true
}

variable "db_username" {
  type    = string
  default = "metabase"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "db_admin_username" {
  type    = string
  default = "admin"
}

variable "db_admin_password" {
  type        = string
  sensitive   = true
  description = "The password for the database administrator. This value is sensitive and should be kept secure."
}

variable "db_username" {
  type    = string
  default = "metabase"
}

variable "db_password" {
  type        = string
  sensitive   = true
  description = "The password for the Metabase database user."
}

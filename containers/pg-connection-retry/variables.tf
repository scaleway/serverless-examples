variable "db_admin_username" {
  description = "Admin username for the RDB instance"
  type        = string
  default     = "admin"
}

variable "db_admin_password" {
  description = "Admin password for the RDB instance"
  type        = string
  sensitive   = true
}

variable "db_username" {
  description = "Username for the application database user"
  type        = string
  default     = "app"
}

variable "db_password" {
  description = "Password for the application database user"
  type        = string
  sensitive   = true
}

variable "registry_image" {
  description = "Full reference of the container image to deploy (e.g. rg.fr-par.scw.cloud/namespace/name:tag)"
  type        = string
}

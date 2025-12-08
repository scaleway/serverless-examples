variable "mongodb_admin_username" {
  description = "The username for the initial MongoDB user"
  type        = string
  default     = "mongoadmin"
}

variable "mongodb_admin_password" {
  description = "The password for the initial MongoDB user"
  type        = string
  sensitive   = true
}

variable "mongodb_username" {
  description = "The username for the application MongoDB user"
  type        = string
  default     = "appuser"
}

variable "mongodb_password" {
  description = "The password for the application MongoDB user"
  type        = string
  sensitive   = true
}

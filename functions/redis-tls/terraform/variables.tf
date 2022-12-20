variable "redis_user" {
  type = string
}

variable "redis_password" {
  type      = string
  sensitive = true
}

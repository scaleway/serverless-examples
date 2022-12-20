variable "zone" {
  type = string
}

variable "region" {
  type = string
}

variable "env" {
  type = string
}

variable "project_id" {
  type        = string
  description = "Your project ID."
}

variable "auth_token" {
    type = string
    description = "scaleway authtentication token used in the function"
}
variable "stack_env" {
  type    = string
  description = ""
}

variable "app_name" {
  type        = string
  default     = "nm-api-service"
  description = "Name of the app"
}

variable "region" {
  type = string
  description = ""
}

variable "pgp_key" {
  type = string
  description = ""
}

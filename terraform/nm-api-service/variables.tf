variable "stack_env" {
  type = string
}

variable "app_name" {
  type        = string
  default     = "nm-api-service"
  description = "Name of the app"
}

variable "app_version" {
  type        = string
  default     = "1.0.5"
  description = "Container image version used to deploy the lambda function"
}

variable "memory" {
  type        = string
  default     = "128"
  description = "Memory in MB to assign to the Lambda."
}

variable "timeout" {
  type        = string
  default     = "900"
  description = "Seconds in which the Lambda should run before timing out."
}

variable "reserved_concurrency" {
  type        = number
  description = "Reserved concurrency guarantees the maximum number of concurrent instances for the function. A value of 0 disables lambda from being triggered and -1 removes any concurrency limitations. "
  default     = -1
}

variable "allow_headers" {
  description = "Allow headers"
  type        = list(string)

  default = [
    "Accept",
    "Authorization",
    "Cookie",
    "Content-Type",
    "Content-Range",
    "X-Amz-Date",
    "X-Amz-Security-Token",
    "X-Api-Key",
    "X-Requested-With"
  ]
}

variable "allow_methods" {
  description = "Allow methods"
  type        = list(string)

  default = [
    "OPTIONS",
    "HEAD",
    "GET",
    "POST",
    "PUT",
    "PATCH",
    "DELETE",
  ]
}

variable "domain_name" {
  description = "The domain name used to contact this API"
  type        = string
}

variable "allow_origin" {
  description = "Allow origin"
  type        = string
  default     = "*"
}

variable "allow_max_age" {
  description = "Allow response caching time"
  type        = string
  default     = "7200"
}

variable "allow_credentials" {
  description = "Allow credentials"
  default     = true
}
variable "base_path" {
  type    = string
  default = "v1"
}

variable "backup_bucket" {
  type    = string
  default = "nm-backup"
}

variable "license_key" {
  type = string
}

variable "authorized_key" {
  type = string
}

variable "bearer_auth" {
  type = string
}

variable "ses_domain" {
  type = string
}
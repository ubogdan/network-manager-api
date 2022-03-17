variable "stack_region" {
  type        = string
  description = "The region where the application will be deployed"
}

variable "stack_env" {
  type        = string
  description = "The environment where the application will be deployed"
}

variable "app_name" {
  type        = string
  default     = "nm-email-service"
  description = "Name of the app"
}

variable "app_version" {
  type        = string
  default     = "0.0.1"
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

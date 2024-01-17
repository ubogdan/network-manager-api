variable "stack_env" {
  type = string
  default = "prod"
}

variable "backup_bucket" {
  type    = string
  default = "nm-backup"
}

variable "download_bucket" {
  type = string
  default = "nm-release-downloads"
}
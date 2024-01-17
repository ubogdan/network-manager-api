provider "aws" {
  region = "eu-central-1"
}
terraform {
  required_version = ">= 0.14"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }

  backend "s3" {
    encrypt = true
    region  = "eu-central-1"
    key     = "nm-api-infra"
    bucket = "ubogdan-terraform-state"
  }
}

## Software releases download bucket
resource "aws_s3_bucket" "downloads" {
  bucket = var.download_bucket
}

resource "aws_s3_bucket_public_access_block" "downloads" {
  bucket = var.backup_bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

## Clients backup's bucket
resource "aws_s3_bucket" "backup" {
  bucket = var.backup_bucket
}

resource "aws_s3_bucket_public_access_block" "backup" {
  bucket = var.backup_bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "backup" {
  bucket = var.backup_bucket
  rule {
    id     = "LifeCycleRule"
    status = "Enabled"
    expiration {
      days = 90
    }
  }
}

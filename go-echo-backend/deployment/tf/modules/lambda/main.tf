provider "aws" {
  profile = var.profile
  region = var.region
}

locals {
  env_ns      = "${var.name}-${var.env}"
}


data "aws_s3_bucket" "storage" {
  bucket        = "${var.storage_s3_bucket}"
}

data "aws_s3_bucket" "cdn" {
  bucket        = "${var.cdn_s3_bucket}"

}

variable "env" {
  type    = string
}

variable "name" {
  type    = string
}

variable "region" {
  type    = string
}

variable "account_id" {
  type = string
}

variable "certificate_arn" {
  type    = string
}


variable "profile" {
  type    = string
}

variable "trending_domain" {
  type = string
}

variable "trending_s3_bucket" {
  type = string
}

variable "storage_domain" {
  type = string
}

variable "storage_s3_bucket" {
  type = string
}

variable "cdn_domain" {
  type = string
}

variable "cdn_s3_bucket" {
  type = string
}

variable "access_control_allow_origins" {
  type = list(string)
}
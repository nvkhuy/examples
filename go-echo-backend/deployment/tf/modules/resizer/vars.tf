variable "env" {
  type = string
}

variable "name" {
  type = string
}

variable "region" {
  type = string
}

variable "account_id" {
  type = string
}

variable "endpoint_path" {
  type = string
}

variable "certificate_arn" {
  type = string
}


variable "profile" {
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


variable "memory_size" {
  type = number
}

variable "media_jwt_secret" {
  type = string
}

variable "concurrent_executions" {
  type    = number
  default = null
}

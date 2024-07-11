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

variable "media_jwt_secret" {
  type = string
}

variable "function_name" {
  type = string
}

variable "image_name" {
  type = string
}

variable "image_tag" {
  type = string
}

variable "memory_size" {
  type = number
}

variable "concurrent_executions" {
  type = number
}


variable "variables" {
  type = map(string)
}

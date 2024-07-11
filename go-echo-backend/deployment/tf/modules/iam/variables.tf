variable "name" {
  type    = string
  description = "Application name"
  default = null
}

variable "env" {
  type    = string
  description = "Application env, it can be dev/prod/beta"
  default = null
}

variable "state_bucket" {
  type = string
  description = "S3 bucket name of storing state"
}


variable "datastore_bucket" {
  type = string
  description = "S3 bucket name of storing env file"
}


variable "storage_bucket" {
  type = string
  description = "S3 bucket name of storing application assets"
}


variable "cdn_bucket" {
  type = string
  description = "S3 bucket name of storing application assets"
}


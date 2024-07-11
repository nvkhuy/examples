variable "certificate_arn" {}

variable "hosted_zone_id" {
  type    = string
  default = ""
}

variable "webapp_config" {
  type = list(object({
    domain      = string
    bucket_name = string
  }))
}

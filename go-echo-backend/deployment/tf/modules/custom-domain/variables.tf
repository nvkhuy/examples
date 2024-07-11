variable "alb_dns_name" {
  type = string
}

variable "alb_zone_id" {
  type = string
}

variable "alb_listener_http_arn" {
  type = string
}

variable "alb_listener_https_arn" {
  type = string
}

variable "domain" {
  type = string
}

variable "route53_id" {
  type = string
}

variable "redirect" {
  type = object({
    host = string
    path = string
  })
}

variable "path_pattern" {
  type = string
  default = null
}
variable "env" {
  type = string
}

variable "name" {
  type = string
}

variable "region" {
  type = string
}


variable "profile" {
  type = string
}

variable "function_name" {
  type = string
}

variable "route" {
  type = object({
    method = string
    path = string
  })
}

variable "certificate_arn" {
  type = string
  default = null
}

variable "domain_name" {
  type = string
  default = null
}

variable "zone_id" {
  type = string
  default = null
}

variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "internal" {
  type = bool
}

variable "vpc_id" {
  type = string
}

variable "subnets" {
  type = list(string)
}

variable "security_groups" {
  type = list(string)
}

variable "listeners" {
  type = map(object({
    listener_port     = number
    listener_protocol = string
  }))
}

variable "listener_port" {
  type = number
}

variable "listener_protocol" {
  type = string
}

variable "target_groups" {
  type = map(object({
    domain      = optional(string)
    port        = number
    protocol    = string
    host_header = list(string)
    priority    = number
    health_check = optional(object({
      matcher             = optional(string)
      path                = optional(string)
      interval            = optional(number)
      timeout             = optional(number)
      unhealthy_threshold = optional(number)
    }))
  }))
}


variable "certificate_arn" {
  type    = string
  default = null
}

variable "hosted_zone_id" {
  type = string
}

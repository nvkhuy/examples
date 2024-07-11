variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "alb_arn" {
  type = string
}

variable "target_groups" {
  type = map(object({
    port              = number
    protocol          = string
    host_header       = list(string)
    priority          = number
    health_check = optional(object({
      matcher             = optional(string)
      path                = optional(string)
      interval            = optional(number)
      timeout             = optional(number)
      unhealthy_threshold = optional(number)
    }))
  }))
}

variable "listeners" {
  type = map(object({
    listener_port     = string
    listener_protocol = string
  }))
}

variable "certificate_arn" {
  type = string
}

variable "whitelisted_cidr_blocks" {
  type = list(string)
}


variable "instance_id" {
  type = string
}

variable cidr_blocks {
  type = list(string)
}

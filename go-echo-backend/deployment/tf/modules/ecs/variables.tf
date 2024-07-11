variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "region" {
  type = string
}

variable "account_id" {
  type = string
}
variable "ecs_task_execution_role_arn" {
  type = string
}

variable "ecs_task_role_arn" {
    type = string
}
variable "vpc_id" {
  type = string
}

variable "service_config" {
  type = any
}


variable "private_subnets" {
  type = list(string)
}


variable "public_alb_security_group" {
  type = any
  default = null
}


variable "public_alb_target_groups" {
  type = map(object({
    arn = string
  }))
  default = null
}

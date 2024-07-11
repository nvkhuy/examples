variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "region" {
  type = string
}

variable "profile" {
  type = string
}
variable "service_config" {
  type = any
}

variable "account_id" {
  type = string
}

variable "target_groups_name_primary" {
  type = any
}

variable "target_groups_name_secondary" {
  type = any
}

variable "http_listener_arns" {
  type = list(string)
}

variable "https_listener_arns" {
  type = list(string)
}

variable "cluster_name" {
  type = string
}

variable "codestar_arn" {
  type = string
}

variable "ecs_task_execution_role_arn" {
  type = string
}

variable "ecs_task_role_arn" {
    type = string
}
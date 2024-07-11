variable "name" {
  type = string
  description = "Application name"
}

variable "env" {
  type = string
  description = "Environment"
}

variable "connector_arn" {
    type = string
    description = "API Gateway ARN or Load balancer ARN"
}
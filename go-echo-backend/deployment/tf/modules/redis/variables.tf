variable "name" {
  type = string
  description = "Application name"
}

variable "env" {
  type = string
  description = "Environment"
}

variable "vpc_id" {
  type = string
  description = "VPC"
}
variable "private_subnets" {
  type = list(string)
}
variable "num_node_groups" {
  type    = number
  default = null
}

variable "replicas_per_node_group" {
  type    = number
  default = null
}


variable "num_cache_clusters" {
   type    = number
    default = null
}

variable "node_type" {
  type = string
}
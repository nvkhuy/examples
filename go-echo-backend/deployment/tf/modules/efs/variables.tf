variable "vpc_id" {
    type = string
}

variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "subnets" {
  type = list(string)
}


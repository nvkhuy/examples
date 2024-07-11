variable "vpc_id" {
  type = string
  default = ""
}
variable "name" {
  type = string
}

variable "env" {
  type = string
}

variable "cidr" {
  type = string
}

variable "availability_zones" {
  type = list(string)
}

variable "private_subnets" {
  type = list(string)
}

variable "public_subnets" {
  type = list(string)
}
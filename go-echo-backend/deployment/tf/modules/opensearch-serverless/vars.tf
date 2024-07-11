variable "vpc_id" {
  type = string
}

variable "cidr_blocks" {
  type = list(string)
}
variable "name" {
  type = string
}

variable "profile" {
  type = string
}

variable "env" {
  type = string
}

variable "collection_name" {
  type = string
}

variable "hosted_zone_name" {
  type = string
}

variable "logs_domain" {
  type = string
}
variable "cloudwatch_logs" {
  type = list(object({
    name           = string
    log_group_name = string
    filter_pattern = string
  }))
}

variable "saml_metadata_document_url" {
  type = string
  default = ""
}

variable "saml_provider_name" {
  type = string
  default = "okta"
}

variable "saml_group" {
  type = string
  default = "opensearch"
}

variable "users" {
  type = list(string)
  default = []
}


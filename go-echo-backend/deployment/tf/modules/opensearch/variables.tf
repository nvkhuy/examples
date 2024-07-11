variable "name" {
  type        = string
  description = "Application name"
}

variable "env" {
  type        = string
  description = "Environment"
}

variable "profile" {
  type = string
}
variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "cidr_blocks" {
  type = list(string)
}

variable "advanced_security_options_enabled" {
  type        = bool
  description = "If advanced security options is enabled."
  default     = false
}

variable "master_user_name" {
  description = "Master username for accessing OpenSerach."
  type        = string
  default     = "admin"
}

variable "inside_vpc" {
  description = "Openserach inside VPC."
  type        = bool
  default     = false
}

variable "master_password" {
  description = "Master password for accessing OpenSearch. If not specified password will be randomly generated. Password will be stored in AWS `System Manager` -> `Parameter Store` "
  type        = string
  default     = ""
}

variable "master_user_arn" {
  description = "Master user ARN for accessing OpenSearch. If this is set, `advanced_security_options_enabled` must be set to true and  `internal_user_database_enabled` should be set to false."
  type        = string
  default     = ""
}

variable "internal_user_database_enabled" {
  type        = bool
  description = "Internal user database enabled. This should be enabled if we want authentication with master username and master password."
  default     = false
}

variable "volume_type" {
  description = "Volume type of ebs storage."
  type        = string
  default     = "gp2"
}

variable "throughput" {
  description = "Specifies the throughput."
  type        = number
  default     = null
}
variable "ebs_enabled" {
  type        = bool
  description = "EBS enabled"
  default     = true
}

variable "ebs_volume_size" {
  type = number
  default     = 10
}

variable "instance_type" {
  type = string
}

variable "instance_count" {
  type = number
}

variable "dedicated_master_enabled" {
  type    = bool
  default = false
}

variable "dedicated_master_count" {
  type    = number
  default = 0
}

variable "dedicated_master_type" {
  type    = string
  default = null
}

variable "zone_awareness_enabled" {
  type    = bool
  default = false
}

variable "engine_version" {
  type = string
}



variable "cloudwatch_logs" {
  type = list(object({
    name           = string
    log_group_name = string
    filter_pattern = string
  }))
}

variable "node_to_node_encryption_enabled" {
  type = bool
  default = true
}

variable "encrypt_at_rest_enabled" {
  type = bool
  default = true
}

variable "cognito_enabled" {
  description = "Cognito authentification enabled for OpenSearch."
  type        = bool
  default     = false
}


variable "mfa_configuration" {
  type        = string
  description = "Multi-Factor Authentication (MFA) configuration for the User Pool"
  default     = "OFF"
}

variable "allow_unauthenticated_identities" {
  type        = bool
  description = "Allow unauthenticated identities on Cognito Identity Pool"
  default     = true
}

variable "role_mapping" {
  type        = any
  description = "Custom role mapping for identity pool role attachment"
  default     = []
}

variable "auto_software_update_enabled" {
  type = bool

  description = "Whether automatic service software updates are enabled for the domain. Defaults to false."
  default     = false
}

variable "custom_endpoint_enabled" {
  description = "If custom endpoint is enabled."
  type        = bool
  default     = false
}

variable "custom_endpoint" {
  description = "Custom endpoint https."
  type        = string
  default     = ""
}

variable "custom_endpoint_certificate_arn" {
  description = "Custom endpoint certificate."
  type        = string
  default     = null
}

variable "tls_security_policy" {
  description = "TLS security policy."
  type        = string
  default     = "Policy-Min-TLS-1-2-2019-07"
}

variable "create_linked_role" {
  type        = bool
  default     = true
  description = "Should linked role be created"
}

variable "aws_service_name_for_linked_role" {
  type        = string
  description = "AWS service name for linked role."
  default     = "opensearchservice.amazonaws.com"
}

variable "zone_id" {
  type        = string
  description = "Route 53 Zone id."
  default     = ""
}

variable "off_peak_window_enabled" {
  type        = bool
  description = "Enabled the off peak update 10 hour update window. All domains created after Feb 16 2023 will have the off_peak_window_options enabled by default."
  default     = null
}

variable "off_peak_window_start_time" {
  type = object({
    hours   = number
    minutes = number
  })

  description = "Time for the 10h update window to begin. If you don't specify a window start time, AWS will default it to 10:00 P.M. local time."
  default     = null
}
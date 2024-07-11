########################################################################################################################
# Application
variable "profile" {
  type        = string
  description = "AWS profile"
}

variable "hosted_zone_name" {
  type = string
}
variable "region" {
  type        = string
  description = "AWS region"
}

variable "name" {
  type        = string
  description = "Application name"
}

variable "env" {
  type        = string
  description = "Environment"
}

variable "state_bucket" {
  type        = string
  description = "State bucket"
}

variable "state_lock_table" {
  type        = string
  description = "State look table"
}

variable "certificate_arn" {
  type        = string
  description = "Certificate ARN"
}

variable "datastore_bucket" {
  type        = string
  description = "Datastore bucket"
}

variable "storage_bucket" {
  type        = string
  description = "Storage bucket"
}

variable "cdn_bucket" {
  type        = string
  description = "CDN bucket"
}

variable "codestar_arn" {
  type        = string
  description = "AWS Github connection"
}
## CDN
variable "access_control_allow_origins" {
  type = list(string)
}

variable "trending_domain" {
  type = string
}

variable "trending_s3_bucket" {
  type = string
}

variable "storage_domain" {
  type = string
}

variable "storage_s3_bucket" {
  type = string
}
variable "cdn_domain" {
  type = string
}


variable "cdn_s3_bucket" {
  type = string
}


variable "certificate_arn_us_east_1" {
  type = string
}

variable "media_jwt_secret" {
  type = string
}

########################################################################################################################

# VPC
variable "cidr" {
  type        = string
  description = "VPC CIDR"
}

variable "availability_zones" {
  type        = list(string)
  description = "Availability zones that the services are running"
}

variable "private_subnets" {
  type        = list(string)
  description = "Private subnets"
}

variable "public_subnets" {
  type        = list(string)
  description = "Public subnets"
}

########################################################################################################################
#ALB
variable "public_alb_config" {
  type = object({
    name = string
    listeners = map(object({
      listener_port     = number
      listener_protocol = string
    }))
    ingress_rules = list(object({
      from_port   = number
      to_port     = number
      protocol    = string
      cidr_blocks = list(string)
    }))
    egress_rules = list(object({
      from_port   = number
      to_port     = number
      protocol    = string
      cidr_blocks = list(string)
    }))
  })
  description = "Public ALB configuration"
}
########################################################################################################################
# ECS
variable "service_config" {
  type = map(object({
    name                               = string
    image_name                         = string
    container_port                     = number
    host_port                          = number
    port_maping_protocol               = string
    cpu                                = number
    memory                             = number
    desired_count                      = number
    deployment_minimum_healthy_percent = number
    deployment_maximum_percent         = number
    health_check_grace_period_seconds  = number
    entrypoint                         = optional(list(string))
    command                            = optional(list(string))
    map_environment                    = optional(map(string))
    environment = optional(list(object({
      name  = string
      value = string
    })))
    environment_files = optional(list(object({
      value = string
      type  = string
    })))
    readonly_root_filesystem = optional(bool)
    mount_points = optional(list(object({
      containerPath = string
      sourceVolume  = string
      readOnly      = bool
    })))
    efs_volume_name = optional(string)
    efs_volume = optional(object({
      file_system_id = string
      root_directory = string
    }))
    alb_target_group = optional(object({
      domain      = optional(string)
      port        = number
      protocol    = string
      host_header = list(string)
      priority    = optional(number)
      health_check = optional(object({
        matcher             = optional(string)
        path                = optional(string)
        interval            = optional(number)
        timeout             = optional(number)
        unhealthy_threshold = optional(number)
      }))
    }))

    auto_scaling = object({
      max_capacity = number
      min_capacity = number
      cpu = object({
        target_value = number
      })
      memory = object({
        target_value = number
      })
    })
  }))
}
########################################################################################################################
# Open Search
variable "cloudwatch_logs" {
  type = list(object({
    name           = string
    log_group_name = string
    filter_pattern = string
  }))
}

########################################################################################################################
# VPN
variable "vpn_config" {
  type = map(object({
    port        = number
    protocol    = string
    host_header = list(string)
    priority    = number
  }))
}

variable "vpn_instance_id" {
  type = string
}

variable "vpn_ip" {
  type = string
}

variable "vpn_domain" {
  type = string
}

variable "vpn_listeners" {
  type = map(object({
    listener_port     = string
    listener_protocol = string
  }))
}


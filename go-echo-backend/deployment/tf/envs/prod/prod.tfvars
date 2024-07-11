#VPC configurations
cidr               = "10.10.0.0/16"
availability_zones = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
public_subnets     = ["10.10.50.0/24", "10.10.51.0/24", "10.10.52.0/24"]
private_subnets    = ["10.10.0.0/24", "10.10.1.0/24", "10.10.2.0/24"]
web_acl_arn        = "arn:aws:wafv2:ap-southeast-1:809144412580:regional/webacl/wafv2-web-acl/d8087826-c81f-4d81-9639-2e050d7cafe0"

hosted_zone_name        = "joininflow.io"
logs_domain             = "logs"
limited_user_emails     = ["thaitanloi365@gmail.com", "doa9595@gmail.com"]
limited_user_group_name = "inflow-prod-logs-limited-group"
master_user_emails      = ["loithai@joininflow.io"]
master_user_group_name  = "inflow-prod-logs-master-group"
cloudwatch_logs = [
  {
    name           = "inflow-prod-backend-logs"
    log_group_name = "inflow-prod-backend-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-prod-consumer-logs"
    log_group_name = "inflow-prod-consumer-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-prod-website-logs"
    log_group_name = "inflow-prod-website-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-prod-resize"
    log_group_name = "/aws/lambda/inflow-prod-resize"
    filter_pattern = ""
  },
  {
    name           = "inflow-prod-rod"
    log_group_name = "/aws/lambda/inflow-prod-rod"
    filter_pattern = ""
  },
  {
    name           = "inflow-prod-ffmpeg"
    log_group_name = "/aws/lambda/inflow-prod-ffmpeg"
    filter_pattern = ""
  },
  {
    name           = "inflow-beta-backend-logs"
    log_group_name = "inflow-beta-backend-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-beta-consumer-logs"
    log_group_name = "inflow-beta-consumer-logs"
    filter_pattern = ""
  },
]

access_control_allow_origins = [
  "http://localhost:3000",
  "http://localhost:3001",
  "http://localhost:3002",
  "https://beta.joininflow.io",
  "https://beta-admin.joininflow.io",
  "https://beta-seller.joininflow.io",
  "https://beta-brand.joininflow.io",
  "https://brand.joininflow.io",
  "https://admin.joininflow.io",
  "https://seller.joininflow.io",
  "https://www.joininflow.io",
  "https://joininflow.io",
  "https://ai.joininflow.io"
]

#Public ALB configurations
public_alb_config = {
  name = "public-alb"
  listeners = {
    "HTTP" = {
      listener_port     = 80
      listener_protocol = "HTTP"

    },
    "HTTPS" = {
      listener_port     = 443
      listener_protocol = "HTTPS"
    }
  }

  ingress_rules = [
    {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
    {
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },

  ]

  egress_rules = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]
}

#Microservices
service_config = {
  "backend" = {
    name                               = "backend"
    image_name                         = "inflow-prod-backend"
    container_port                     = 8080
    host_port                          = 8080
    port_maping_protocol               = "tcp"
    cpu                                = 512
    memory                             = 1024
    desired_count                      = 1
    deployment_minimum_healthy_percent = 100
    deployment_maximum_percent         = 200
    health_check_grace_period_seconds  = 10
    command                            = null
    entrypoint                         = null
    environment                        = null
    environment_files                  = null
    mount_points = [{
      containerPath = "/app/efs",
      sourceVolume  = "inflow-prod-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-prod-efs"
    efs_volume = {
      file_system_id = "fs-0cbe4b928fcf2d6c6"
      root_directory = "/"
    }
    alb_target_group = {
      port              = 8080
      protocol          = "HTTP"
      health_check_path = "/health_check"
      priority          = 1
      host_header       = ["api.joininflow.io"]
      domain            = "api.joininflow.io"
    }
    auto_scaling = {
      max_capacity = 4
      min_capacity = 1
      cpu = {
        target_value = 75
      }
      memory = {
        target_value = 75
      }
    }
  },
  "consumer" = {
    name                               = "consumer"
    image_name                         = "inflow-prod-consumer"
    container_port                     = 8080
    host_port                          = 8080
    port_maping_protocol               = "tcp"
    cpu                                = 512
    memory                             = 1024
    desired_count                      = 1
    deployment_minimum_healthy_percent = 100
    deployment_maximum_percent         = 200
    health_check_grace_period_seconds  = 10
    command                            = null
    entrypoint                         = null
    environment                        = null
    environment_files                  = null
    mount_points = [{
      containerPath = "/app/efs",
      sourceVolume  = "inflow-prod-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-prod-efs"
    efs_volume = {
      file_system_id = "fs-0cbe4b928fcf2d6c6"
      root_directory = "/"
    }
    alb_target_group = {
      port              = 8080
      protocol          = "HTTP"
      host_header       = ["consumer.joininflow.io"]
      health_check_path = "/health_check"
      priority          = 1
      domain            = "consumer.joininflow.io"
    }
    auto_scaling = {
      max_capacity = 4
      min_capacity = 1
      cpu = {
        target_value = 75
      }
      memory = {
        target_value = 75
      }
    }
  }
}

# VPN
vpn_instance_id = "i-09a267c6ccdc22229"
vpn_ip          = "52.76.40.168"
vpn_domain      = "vpn.joininflow.io"
vpn_config = {
  "vpn-http" = {
    port              = 80
    protocol          = "HTTP"
    health_check_path = "/heath_check"
    priority          = 1
    host_header       = ["vpn.joininflow.io"]
  },

}

vpn_listeners = {
  "HTTP" = {
    listener_port     = 80
    listener_protocol = "HTTP"

  },
  "HTTPS" = {
    listener_port     = 443
    listener_protocol = "HTTPS"
  }
}

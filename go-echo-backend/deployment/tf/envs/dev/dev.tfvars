# VPC configurations
codestar_arn       = "arn:aws:codestar-connections:ap-southeast-1:809144412580:connection/de7f3e1e-579c-45f8-bc89-a96fdbfc52a3"
cidr               = "10.10.0.0/16"
availability_zones = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
public_subnets     = ["10.10.50.0/24", "10.10.51.0/24", "10.10.52.0/24"]
private_subnets    = ["10.10.0.0/24", "10.10.1.0/24", "10.10.2.0/24"]
hosted_zone_name   = "joininflow.io"

# Open Search
logs_domain             = "dev-logs"
limited_user_emails     = ["thaitanloi365@gmail.com", "doa9595@gmail.com"]
limited_user_group_name = "inflow-dev-logs-limited-group"
master_user_emails      = ["loithai@joininflow.io", "huynguyen@joininflow.io"]
master_user_group_name  = "inflow-dev-logs-master-group"
cloudwatch_logs = [
  {
    name           = "inflow-dev-backend-logs"
    log_group_name = "inflow-dev-backend-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-consumer-logs"
    log_group_name = "inflow-dev-consumer-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-crawler-logs"
    log_group_name = "inflow-dev-crawler-logs"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-resize"
    log_group_name = "/aws/lambda/inflow-dev-resize"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-rod"
    log_group_name = "/aws/lambda/inflow-dev-rod"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-ffmpeg"
    log_group_name = "/aws/lambda/inflow-dev-ffmpeg"
    filter_pattern = ""
  },
  {
    name           = "inflow-dev-waf"
    log_group_name = "aws-waf-logs-wafv2-web-acl"
    filter_pattern = ""
  },
]

access_control_allow_origins = [
  "http://localhost:3000",
  "http://localhost:3001",
  "http://localhost:3002",
  "https://dev.joininflow.io",
  "https://dev-admin.joininflow.io",
  "https://dev-seller.joininflow.io",
  "https://dev-brand.joininflow.io",
  "https://dev-ai.joininflow.io",
  "https://dev-crawler.joininflow.io"
]
# Public ALB configurations
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

# Microservices
service_config = {
  "backend" = {
    name                               = "backend"
    image_name                         = "inflow-dev-backend"
    container_port                     = 8080
    host_port                          = 8080
    port_maping_protocol               = "tcp"
    cpu                                = 256
    memory                             = 512
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
      sourceVolume  = "inflow-dev-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-dev-efs"
    efs_volume = {
      file_system_id = "fs-0744628a84a9edccd"
      root_directory = "/"
    }
    alb_target_group = {
      domain      = "dev-api.joininflow.io"
      port        = 8080
      protocol    = "HTTP"
      priority    = 1
      host_header = ["dev-api.joininflow.io"],
      health_check = {
        matcher             = "200,301,302"
        path                = "/"
        interval            = 300
        timeout             = 120
        unhealthy_threshold = 5
      }
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
    image_name                         = "inflow-dev-consumer"
    container_port                     = 8080
    host_port                          = 8080
    port_maping_protocol               = "tcp"
    cpu                                = 256
    memory                             = 512
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
      sourceVolume  = "inflow-dev-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-dev-efs"
    efs_volume = {
      file_system_id = "fs-0744628a84a9edccd"
      root_directory = "/"
    }
    alb_target_group = {
      domain      = "dev-consumer.joininflow.io"
      port        = 8080
      protocol    = "HTTP"
      host_header = ["dev-consumer.joininflow.io"]
      priority    = 2
      health_check = {
        matcher             = "200,301,302"
        path                = "/"
        interval            = 300
        timeout             = 120
        unhealthy_threshold = 5
      }
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
  ## Crawler
  "crawler" = {
    name                               = "crawler"
    image_name                         = "inflow-dev-crawler"
    container_port                     = 8080
    host_port                          = 8080
    port_maping_protocol               = "tcp"
    cpu                                = 1024
    memory                             = 2048
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
      sourceVolume  = "inflow-dev-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-dev-efs"
    efs_volume = {
      file_system_id = "fs-0744628a84a9edccd"
      root_directory = "/"
    }
    alb_target_group = {
      domain      = "dev-crawler.joininflow.io"
      port        = 8080
      protocol    = "HTTP"
      host_header = ["dev-crawler.joininflow.io"]
      priority    = 2
      health_check = {
        matcher             = "200,301,302"
        path                = "/"
        interval            = 300
        timeout             = 120
        unhealthy_threshold = 5
      }
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
}

# VPN
vpn_instance_id = "i-0f9acae7b37b42d60"
vpn_domain      = "dev-vpn.joininflow.io"
vpn_ip          = "18.139.215.157"
vpn_config = {
  "vpn-http" = {
    port              = 80
    protocol          = "HTTP"
    health_check_path = "/"
    priority          = 1
    host_header       = ["dev-vpn.joininflow.io"]
  }
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

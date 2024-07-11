#VPC configurations
vpc_id             = "vpc-0970f7741a9e391d8"
availability_zones = ["ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"]
hosted_zone_name   = "joininflow.io"
web_acl_arn = "arn:aws:wafv2:ap-southeast-1:809144412580:regional/webacl/wafv2-web-acl/d8087826-c81f-4d81-9639-2e050d7cafe0"

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
    image_name                         = "inflow-beta-backend"
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
      sourceVolume  = "inflow-beta-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-beta-efs"
    efs_volume = {
      file_system_id = "fs-0cbe4b928fcf2d6c6"
      root_directory = "/"
    }
    alb_target_group = {
      port              = 8080
      protocol          = "HTTP"
      health_check_path = "/health_check"
      priority          = 1
      domain            = "beta-api.joininflow.io"
      host_header       = ["beta-api.joininflow.io"]
    }
    auto_scaling = {
      max_capacity = 2
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
    image_name                         = "inflow-beta-consumer"
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
      sourceVolume  = "inflow-beta-efs"
      readOnly      = false
    }]
    efs_volume_name = "inflow-beta-efs"
    efs_volume = {
      file_system_id = "fs-0cbe4b928fcf2d6c6"
      root_directory = "/"
    }
    alb_target_group = {
      port              = 8080
      protocol          = "HTTP"
      domain            = "beta-consumer.joininflow.io"
      host_header       = ["beta-consumer.joininflow.io"]
      health_check_path = "/health_check"
      priority          = 1
    }
    auto_scaling = {
      max_capacity = 2
      min_capacity = 1
      cpu = {
        target_value = 75
      }
      memory = {
        target_value = 75
      }
    }
  },
  # ## Websites
  # "website" = {
  #   name                               = "website"
  #   image_name                         = "inflow-beta-website"
  #   container_port                     = 3000
  #   host_port                          = 3000
  #   port_maping_protocol               = "tcp"
  #   cpu                                = 256
  #   memory                             = 512
  #   desired_count                      = 1
  #   deployment_minimum_healthy_percent = 100
  #   deployment_maximum_percent         = 200
  #   health_check_grace_period_seconds  = 10
  #   command                            = null
  #   entrypoint                         = null
  #   environment                        = null
  #   environment_files                  = null
  #   readonly_root_filesystem           = false
  #   mount_points                       = null
  #   efs_volume_name                    = null
  #   efs_volume                         = null
  #   alb_target_group = {
  #     domain      = "beta.joininflow.io"
  #     port        = 80
  #     protocol    = "HTTP"
  #     priority    = 1
  #     host_header = ["beta.joininflow.io"]
  #     health_check = {
  #       matcher             = "200,301,302"
  #       path                = "/"
  #       interval            = 300
  #       timeout             = 120
  #       unhealthy_threshold = 5
  #     }
  #   }
  #   auto_scaling = {
  #     max_capacity = 2
  #     min_capacity = 1
  #     cpu = {
  #       target_value = 75
  #     }
  #     memory = {
  #       target_value = 75
  #     }
  #   }
  # },
  # "brand" = {
  #   name                               = "brand"
  #   image_name                         = "inflow-beta-brand"
  #   container_port                     = 3000
  #   host_port                          = 3000
  #   port_maping_protocol               = "tcp"
  #   cpu                                = 256
  #   memory                             = 512
  #   desired_count                      = 1
  #   deployment_minimum_healthy_percent = 100
  #   deployment_maximum_percent         = 200
  #   health_check_grace_period_seconds  = 10
  #   command                            = null
  #   entrypoint                         = null
  #   environment                        = null
  #   environment_files                  = null
  #   readonly_root_filesystem           = false
  #   mount_points                       = null
  #   efs_volume_name                    = null
  #   efs_volume                         = null
  #   alb_target_group = {
  #     domain      = "beta-brand.joininflow.io"
  #     port        = 80
  #     protocol    = "HTTP"
  #     priority    = 1
  #     host_header = ["beta-brand.joininflow.io"]
  #     health_check = {
  #       matcher             = "200,301,302"
  #       path                = "/"
  #       interval            = 300
  #       timeout             = 120
  #       unhealthy_threshold = 5
  #     }
  #   }
  #   auto_scaling = {
  #     max_capacity = 2
  #     min_capacity = 1
  #     cpu = {
  #       target_value = 75
  #     }
  #     memory = {
  #       target_value = 75
  #     }
  #   }
  # }
  # "admin" = {
  #   name                               = "admin"
  #   image_name                         = "inflow-beta-admin"
  #   container_port                     = 3000
  #   host_port                          = 3000
  #   port_maping_protocol               = "tcp"
  #   cpu                                = 256
  #   memory                             = 512
  #   desired_count                      = 1
  #   deployment_minimum_healthy_percent = 100
  #   deployment_maximum_percent         = 200
  #   health_check_grace_period_seconds  = 10
  #   command                            = null
  #   entrypoint                         = null
  #   environment                        = null
  #   environment_files                  = null
  #   readonly_root_filesystem           = false
  #   mount_points                       = null
  #   efs_volume_name                    = null
  #   efs_volume                         = null
  #   alb_target_group = {
  #     port        = 80
  #     protocol    = "HTTP"
  #     priority    = 1
  #     host_header = ["beta-admin.joininflow.io"]
  #     domain      = "beta-admin.joininflow.io"
  #     health_check = {
  #       matcher             = "200,301,302"
  #       path                = "/"
  #       interval            = 300
  #       timeout             = 120
  #       unhealthy_threshold = 5
  #     }
  #   }
  #   auto_scaling = {
  #     max_capacity = 2
  #     min_capacity = 1
  #     cpu = {
  #       target_value = 75
  #     }
  #     memory = {
  #       target_value = 75
  #     }
  #   }
  # },
  # "seller" = {
  #   name                               = "seller"
  #   image_name                         = "inflow-beta-seller"
  #   container_port                     = 3000
  #   host_port                          = 3000
  #   port_maping_protocol               = "tcp"
  #   cpu                                = 256
  #   memory                             = 512
  #   desired_count                      = 1
  #   deployment_minimum_healthy_percent = 100
  #   deployment_maximum_percent         = 200
  #   health_check_grace_period_seconds  = 10
  #   command                            = null
  #   entrypoint                         = null
  #   environment                        = null
  #   environment_files                  = null
  #   readonly_root_filesystem           = false
  #   mount_points                       = null
  #   efs_volume_name                    = null
  #   efs_volume                         = null
  #   alb_target_group = {
  #     port        = 80
  #     protocol    = "HTTP"
  #     priority    = 1
  #     host_header = ["beta-seller.joininflow.io"]
  #     domain      = "beta-seller.joininflow.io"
  #     health_check = {
  #       matcher             = "200,301,302"
  #       path                = "/"
  #       interval            = 300
  #       timeout             = 120
  #       unhealthy_threshold = 5
  #     }
  #   }
  #   auto_scaling = {
  #     max_capacity = 2
  #     min_capacity = 1
  #     cpu = {
  #       target_value = 75
  #     }
  #     memory = {
  #       target_value = 75
  #     }
  #   }
  # }
}


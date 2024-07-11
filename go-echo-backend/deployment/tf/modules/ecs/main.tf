locals {
  env_ns       = "${var.name}-${var.env}"
  app_services = [for service, config in var.service_config : service]

}

resource "aws_ecs_cluster" "ecs_cluster" {
  name = lower("${local.env_ns}-cluster")
  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_cloudwatch_log_group" "ecs_cloudwatch_log_group" {
  for_each          = toset(local.app_services)
  name              = lower("${local.env_ns}-${each.key}-logs")
  retention_in_days = 14
}

resource "time_static" "this" {}

module "container_definition" {
  for_each = var.service_config

  source = "cloudposse/ecs-container-definition/aws"
  # Cloud Posse recommends pinning every module to a specific version
  version = "0.61.1"

  container_name   = each.value.name
  container_image  = "${var.account_id}.dkr.ecr.${var.region}.amazonaws.com/${each.value.image_name}:latest"
  essential        = true
  container_cpu    = each.value.cpu
  container_memory = each.value.memory

  environment_files        = each.value.environment_files
  environment              = each.value.environment
  entrypoint               = each.value.entrypoint
  command                  = each.value.command
  readonly_root_filesystem = each.value.readonly_root_filesystem

  port_mappings = [
    {
      containerPort = each.value.container_port
      hostPort      = each.value.host_port
      protocol      = each.value.port_maping_protocol
    }
  ]

  mount_points = each.value.mount_points

  log_configuration = {
    logDriver = "awslogs"
    options = {
      awslogs-group         = "${local.env_ns}-${each.value["name"]}-logs"
      awslogs-region        = var.region
      awslogs-stream-prefix = local.env_ns
    }
  }

  map_environment = merge({
    _key        = each.key
    _build_time = time_static.this.rfc3339
  }, each.value.map_environment)
}

#Create task definitions for app services
resource "aws_ecs_task_definition" "ecs_task_definition" {
  for_each                 = var.service_config
  family                   = "${local.env_ns}-${each.key}"
  task_role_arn            = var.ecs_task_role_arn
  execution_role_arn       = var.ecs_task_execution_role_arn
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  memory                   = each.value.memory
  cpu                      = each.value.cpu

  container_definitions = module.container_definition[each.key].json_map_encoded_list

  dynamic "volume" {
    for_each = each.value.efs_volume == null ? [] : [true]
    content {
      name = each.value.efs_volume_name
      efs_volume_configuration {
        file_system_id = each.value.efs_volume.file_system_id
        root_directory = each.value.efs_volume.root_directory
      }
    }
  }
}


#Create services for app services
resource "aws_ecs_service" "private_service" {
  for_each = var.service_config

  name                               = "${local.env_ns}-${each.value.name}"
  cluster                            = aws_ecs_cluster.ecs_cluster.id
  task_definition                    = aws_ecs_task_definition.ecs_task_definition[each.key].arn
  launch_type                        = "FARGATE"
  desired_count                      = each.value.desired_count
  deployment_minimum_healthy_percent = each.value.deployment_minimum_healthy_percent
  deployment_maximum_percent         = each.value.deployment_maximum_percent
  health_check_grace_period_seconds  = each.value.health_check_grace_period_seconds

  network_configuration {
    subnets          = var.private_subnets
    assign_public_ip = false
    security_groups = [
      aws_security_group.service_security_group.id
    ]
  }

  load_balancer {
    target_group_arn = var.public_alb_target_groups[each.key].arn
    container_name   = each.value.name
    container_port   = each.value.container_port
  }

}

resource "aws_appautoscaling_target" "service_autoscaling" {
  for_each           = var.service_config
  max_capacity       = each.value.auto_scaling.max_capacity
  min_capacity       = each.value.auto_scaling.min_capacity
  resource_id        = "service/${aws_ecs_cluster.ecs_cluster.name}/${aws_ecs_service.private_service[each.key].name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "ecs_policy_memory" {
  for_each           = var.service_config
  name               = "${local.env_ns}-memory-autoscaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.service_autoscaling[each.key].resource_id
  scalable_dimension = aws_appautoscaling_target.service_autoscaling[each.key].scalable_dimension
  service_namespace  = aws_appautoscaling_target.service_autoscaling[each.key].service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }

    target_value = each.value.auto_scaling.memory.target_value
    scale_in_cooldown = 300
    scale_out_cooldown = 300
  }
}

resource "aws_appautoscaling_policy" "ecs_policy_cpu" {
  for_each           = var.service_config
  name               = "${local.env_ns}-cpu-autoscaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.service_autoscaling[each.key].resource_id
  scalable_dimension = aws_appautoscaling_target.service_autoscaling[each.key].scalable_dimension
  service_namespace  = aws_appautoscaling_target.service_autoscaling[each.key].service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }

    target_value = each.value.auto_scaling.cpu.target_value
    scale_in_cooldown = 300
    scale_out_cooldown = 300
  }
}

resource "aws_security_group" "service_security_group" {
  vpc_id = var.vpc_id

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [var.public_alb_security_group.security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


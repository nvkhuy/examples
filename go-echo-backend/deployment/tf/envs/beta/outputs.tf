output "public_alb_target_groups" {
 value = module.public_alb.target_groups
}

output "aws_cloudwatch_log" {
 value = module.ecs.aws_cloudwatch_log_group
}


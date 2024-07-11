output "vpc_id" {
 value = module.vpc.vpc_id
}

output "private_subnets" {
 value = module.vpc.private_subnets
}

output "public_subnets" {
 value = module.vpc.public_subnets
}

output "public_alb_target_groups" {
 value = module.public_alb.target_groups
}

output "aws_cloudwatch_log" {
 value = module.ecs.aws_cloudwatch_log_group
}


output "redis" {
 value = module.redis
}

output "efs" {
 value = module.efs
}
output "cdn" {
 value = module.cdn
}

output "vpn" {
 value = module.vpn
}

output "opensearch" {
 value = module.opensearch
}

output "apigateway_rod" {
 value = module.apigateway_rod
}

output "apigateway_resize" {
 value = module.apigateway_resize
}

output "apigateway_ffmpeg" {
 value = module.apigateway_ffmpeg
}


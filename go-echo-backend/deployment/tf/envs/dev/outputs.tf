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

output "public_alb_listener_http" {
 value = module.public_alb.aws_alb_listener_http
}

output "public_alb_listener_https" {
 value = module.public_alb.aws_alb_listener_https
}
output "aws_cloudwatch_log" {
 value = module.ecs.aws_cloudwatch_log_group
}

output "redis_primary_endpoint_address" {
 value = module.redis.redis_primary_endpoint_address
}

output "redis_crawler_primary_endpoint_address" {
 value = module.redis_crawler.redis_primary_endpoint_address
}

output "efs" {
 value = module.efs
}

output "cdn" {
 value = module.cdn
}


output "opensearch" {
 value = module.opensearch
}

output "waf" {
 value = module.waf
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

output "apigateway_blur" {
 value = module.apigateway_blur
}
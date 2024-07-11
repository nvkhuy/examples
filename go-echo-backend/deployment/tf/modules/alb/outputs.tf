output "alb_dns_name" {
  value = aws_alb.alb.dns_name
}
output "alb_zone_id" {
  value = aws_alb.alb.zone_id
}

output "alb_id" {
  value = aws_alb.alb.id
}

output "alb_arn" {
  value = aws_alb.alb.arn
}

output "target_groups" {
  value = aws_alb_target_group.alb_target_group
}

# output "target_groups_secondary" {
#   value = aws_alb_target_group.alb_target_group_secondary
# }

output "aws_alb_listener_http" {
  value = aws_alb_listener.alb_listener["HTTP"].arn
}

output "aws_alb_listener_https" {
  value = aws_alb_listener.alb_listener["HTTPS"].arn
}
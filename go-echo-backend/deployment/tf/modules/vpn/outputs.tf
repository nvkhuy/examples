output "target_groups" {
  value = aws_alb_target_group.alb_target_group
}

output "aws_alb_listener" {
  value = aws_alb_listener.alb_listener
}


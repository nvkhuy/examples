
output "ecs_exec_role_arn" {
  value = aws_iam_role.ecs_exec.arn
}

output "ecs_task_arn" {
  value = aws_iam_role.ecs_task.arn
}

output "ecs_exec_instance_profile_name" {
  value = aws_iam_instance_profile.ecs_exec.name
}

output "opensearch_endpoint" {
  value = aws_opensearch_domain.opensearch.endpoint
}


output "master_user_role_arn" {
  value = aws_iam_role.master_user_role.arn
}
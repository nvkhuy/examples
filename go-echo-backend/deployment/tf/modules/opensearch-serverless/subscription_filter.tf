resource "aws_cloudwatch_log_subscription_filter" "default" {
  count           = length(var.cloudwatch_logs)
  name            = var.cloudwatch_logs[count.index].name
  log_group_name  = var.cloudwatch_logs[count.index].log_group_name
  filter_pattern  = var.cloudwatch_logs[count.index].filter_pattern
  destination_arn = aws_lambda_function.central_logs_lambda.arn
  depends_on      = [aws_lambda_function.central_logs_lambda]
}

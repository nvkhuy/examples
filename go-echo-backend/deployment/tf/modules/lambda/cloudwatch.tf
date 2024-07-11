resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.func.function_name}"
  retention_in_days = 30
  lifecycle {
    prevent_destroy = false
  }
}

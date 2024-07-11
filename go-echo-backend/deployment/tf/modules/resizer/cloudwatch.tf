resource "aws_cloudwatch_log_group" "resizer" {
  name              = "/aws/lambda/${var.name}-${var.env}-resizer"
  retention_in_days = 30
  lifecycle {
    prevent_destroy = false
  }
}

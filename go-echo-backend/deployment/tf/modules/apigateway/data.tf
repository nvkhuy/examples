data "aws_lambda_function" "func" {
  function_name = "${var.name}-${var.env}-${var.function_name}"
}

data "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${data.aws_lambda_function.func.function_name}"
}
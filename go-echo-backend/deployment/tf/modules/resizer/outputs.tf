
output "storage_bucket" {
  value = data.aws_s3_bucket.storage
}

output "cdn_bucket" {
  value = data.aws_s3_bucket.cdn
}

output "function_resizer_arn" {
  value = "${aws_lambda_function.resizer.arn}:${aws_lambda_function.resizer.version}"
}

output "function_resizer_name" {
  value = aws_lambda_function.resizer.function_name
}

output "api_gateway_integration_resizer" {
  value = aws_apigatewayv2_integration.resizer
}

output "api_gateway_rest_api_resizer" {
  value = aws_apigatewayv2_api.resizer
}

output "endpoint_url" {
  value = "${aws_apigatewayv2_stage.resizer.invoke_url}/${var.endpoint_path}?key=64w/test.jpeg"
}

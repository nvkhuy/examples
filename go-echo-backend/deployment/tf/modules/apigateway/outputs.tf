output "endpoint_url" {
  value = "${aws_apigatewayv2_stage.stage.invoke_url}${var.route.path}"
}

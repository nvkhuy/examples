resource "aws_apigatewayv2_api" "resizer" {
  name          = "${var.name}_${var.env}_resizer"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "resizer" {
  api_id                 = aws_apigatewayv2_api.resizer.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.resizer.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "resizer" {
  api_id    = aws_apigatewayv2_api.resizer.id
  route_key = "GET /${var.endpoint_path}"
  target    = "integrations/${aws_apigatewayv2_integration.resizer.id}"
}

resource "aws_apigatewayv2_stage" "resizer" {
  api_id      = aws_apigatewayv2_api.resizer.id
  name        = var.env
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.resizer.arn
    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )

  }
}

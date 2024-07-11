resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "${data.aws_lambda_function.func.function_name}-AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.func.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_api" "api" {
  name          = "${var.name}-${var.env}-${var.function_name}"
  protocol_type = "HTTP"

}

resource "aws_apigatewayv2_integration" "integration" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = data.aws_lambda_function.func.invoke_arn
  payload_format_version = "2.0"

}

resource "aws_apigatewayv2_route" "route" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "${var.route.method} ${var.route.path}"
  target    = "integrations/${aws_apigatewayv2_integration.integration.id}"
}


resource "aws_apigatewayv2_stage" "stage" {
  api_id    = aws_apigatewayv2_api.api.id

  name        = var.env
  auto_deploy = true
  access_log_settings {
    destination_arn = data.aws_cloudwatch_log_group.log_group.arn
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

resource "aws_apigatewayv2_domain_name" "api" {
  count = var.domain_name != null && var.certificate_arn != null ? 1 : 0
  domain_name = var.domain_name
  domain_name_configuration {
    certificate_arn = var.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}


resource "aws_apigatewayv2_api_mapping" "api" {
  count = var.domain_name != null && var.certificate_arn != null ? 1 : 0

  api_id      = aws_apigatewayv2_api.api.id
  domain_name = aws_apigatewayv2_domain_name.api[0].id
  stage       = aws_apigatewayv2_stage.stage.id
}




resource "aws_route53_record" "domain_record" {
  count = var.domain_name != null && var.zone_id != null ? 1 : 0

  name    = aws_apigatewayv2_domain_name.api[0].domain_name
  type    = "A"
  zone_id = var.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
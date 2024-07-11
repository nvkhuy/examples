# output "endpoint_url" {
#   value = [ for item in var.functions: "${lookup(aws_apigatewayv2_stage.stage,item.name).invoke_url}${lookup(var.functions,item.name).endpoint_path}"]
# }

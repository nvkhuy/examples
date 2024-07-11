resource "aws_security_group" "lambda" {
  name_prefix = "${var.name}-${var.env}-lambda-sg"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = var.cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = var.cidr_blocks
  }
}

resource "aws_lambda_permission" "central_logs_log_group_invoke_lambda" {
  statement_id   = "${local.namespace}-lambda-permissions"
  action         = "lambda:InvokeFunction"
  function_name  = aws_lambda_function.central_logs_lambda.function_name
  principal      = "logs.${data.aws_region.current.name}.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id

  depends_on = [aws_lambda_function.central_logs_lambda]
}




resource "aws_lambda_function" "central_logs_lambda" {
  function_name    = local.function_name
  role             = aws_iam_role.master_user_role.arn
  publish          = true
  timeout          = 60
  memory_size      = 256
  package_type     = "Image"
  image_uri        = "${local.repo_uri}/${local.image_name}:${local.image_tag}"
  source_code_hash = trimprefix(data.aws_ecr_repository.application_ecr_repo.id, "sha256:")
  
  dynamic "vpc_config" {
    for_each = var.inside_vpc ? [1] : []
    content {
      subnet_ids         = local.subnet_ids
      security_group_ids = [aws_security_group.es[0].id]
    }
  }

  environment {
    variables = {
      ENDPOINT          = aws_opensearch_domain.opensearch.endpoint
      SERVICE           = "es" // "aoss" for Amazon OpenSearch Serverless
      INDEX_NAME_PREFIX = "cwl"
    }
  }


  depends_on = [
    null_resource.ecr_image
  ]
}

# resource "aws_cloudwatch_log_group" "function_log_group" {
#   name              = "/aws/lambda/${aws_lambda_function.central_logs_lambda.function_name}"
#   retention_in_days = 7
#   lifecycle {
#     prevent_destroy = false
#   }
# }

resource "null_resource" "ecr_image" {
  triggers = {
    dir_sha1 = sha1(join("", [for f in fileset("${path.module}/lambda", "*"): filesha1("${path.module}/lambda/${f}")]))
  }

  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --profile ${var.profile} --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${local.repo_uri}
           cd ${path.module}/lambda
           docker build --no-cache --platform linux/amd64 -t ${data.aws_ecr_repository.application_ecr_repo.repository_url}:${local.image_tag} .
           docker push ${data.aws_ecr_repository.application_ecr_repo.repository_url}:${local.image_tag}
       EOF
  }
}

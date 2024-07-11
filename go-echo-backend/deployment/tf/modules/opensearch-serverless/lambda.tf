resource "aws_security_group" "lambda" {
  name_prefix = "${var.name}-${var.env}-lambda-opensearch-serverless"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


resource "aws_lambda_permission" "central_logs_log_group_invoke_lambda" {
  statement_id   = "${local.namespace}-lambda-permissions"
  action         = "lambda:InvokeFunction"
  function_name  = aws_lambda_function.central_logs_lambda.function_name
  principal      = "logs.${data.aws_region.current.name}.amazonaws.com"
  source_account = data.aws_caller_identity.current.account_id
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

  tracing_config {
    mode = "Active" # Activate AWS X-Ray
  }

  vpc_config {
    subnet_ids         = data.aws_subnets.private.ids
    security_group_ids = [aws_security_group.lambda.id]
  }

  environment {
    variables = {
      ENDPOINT   = aws_opensearchserverless_collection.collection.collection_endpoint
      INDEX_NAME = "cwl"
    }
  }


  depends_on = [
    null_resource.ecr_image,
  ]
}

resource "null_resource" "ecr_image" {
  triggers = {
    timestamp = timestamp()
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

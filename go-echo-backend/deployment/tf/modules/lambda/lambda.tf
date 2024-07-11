resource "aws_lambda_function" "func" {
  function_name                  = "${local.env_ns}-${var.function_name}"
  role                           = aws_iam_role.lambda.arn
  publish                        = true
  timeout                        = 60
  memory_size                    = var.memory_size
  package_type                   = "Image"
  image_uri                      = "${var.account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.image_name}:${var.image_tag}"
  reserved_concurrent_executions = var.concurrent_executions
  source_code_hash               = trimprefix(data.aws_ecr_repository.application_ecr_repo.id, "sha256:")
  
  tracing_config {
    mode = "Active" # Activate AWS X-Ray
  }

  environment {
    variables = var.variables
  }


  depends_on = [
    aws_iam_role_policy_attachment.lambda,
    null_resource.ecr_image,
  ]
}


resource "null_resource" "ecr_image" {
  triggers = {
    timestamp   = timestamp()
  }

  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --profile ${var.profile} --region ${var.region} | docker login --username AWS --password-stdin ${var.account_id}.dkr.ecr.${var.region}.amazonaws.com
           cd ${path.module}/${var.function_name}
           docker build --no-cache --platform linux/amd64 -t ${data.aws_ecr_repository.application_ecr_repo.repository_url}:${var.image_tag} .
           docker push ${data.aws_ecr_repository.application_ecr_repo.repository_url}:${var.image_tag}
       EOF
  }
}



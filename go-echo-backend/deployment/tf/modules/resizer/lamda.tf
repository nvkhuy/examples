resource "null_resource" "build" {
  provisioner "local-exec" {
    command = "cd ${path.module}/handler && npm i"
  }
}

data "archive_file" "resizer_file" {
  type        = "zip"
  source_dir  = "${path.module}/handler"
  output_path = "${path.module}/handler.zip"
  depends_on  = [null_resource.build]
}


resource "aws_s3_bucket_object" "layer_upload" {
  bucket     = data.aws_s3_bucket.cdn.bucket
  key        = "ffmpeg.zip"
  source     = "${path.module}/layers/ffmpeg.zip"
  etag       = filebase64sha256("${path.module}/layers/ffmpeg.zip")
  depends_on = [null_resource.build]
}

resource "aws_lambda_layer_version" "ffmpeg" {
  layer_name = "ffmpeg"
  s3_bucket  = data.aws_s3_bucket.cdn.bucket
  s3_key     = aws_s3_bucket_object.layer_upload.key

  source_code_hash    = filebase64sha256("${path.module}/layers/ffmpeg.zip")
  compatible_runtimes = ["nodejs16.x"]

  depends_on = [aws_s3_bucket_object.layer_upload]
}

resource "aws_lambda_function" "resizer" {
  filename         = data.archive_file.resizer_file.output_path
  function_name    = "${var.name}-${var.env}-resizer"
  role             = aws_iam_role.lambda.arn
  handler          = "index.handler"
  runtime          = "nodejs16.x"
  source_code_hash = data.archive_file.resizer_file.output_base64sha256
  publish          = true
  timeout          = 60
  memory_size      = var.memory_size

  reserved_concurrent_executions = var.concurrent_executions
  tracing_config {
    mode = "Active" # Activate AWS X-Ray
  }


  layers = [aws_lambda_layer_version.ffmpeg.arn]

  environment {
    variables = {
      AWS_ORIGIN_BUCKET = var.storage_s3_bucket
      AWS_ORIGIN_REGION = var.region
      AWS_DEST_BUCKET   = var.cdn_s3_bucket
      AWS_DEST_REGION   = var.region
      AWS_DEST_URL      = "https://${var.cdn_domain}"
      AWS_STORAGE_URL   = "https://${var.storage_domain}"
      JWT_SECRET        = var.media_jwt_secret
    }
  }


  depends_on = [
    aws_iam_role_policy_attachment.lambda,
    aws_cloudwatch_log_group.resizer,
  ]
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.resizer.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.resizer.execution_arn}/*/*"
}

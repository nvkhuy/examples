resource "aws_iam_role" "lambda" {
  name = "${local.env_ns_full}-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "lambda.amazonaws.com",
          "edgelambda.amazonaws.com"
        ]
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda" {
  name = "${local.env_ns_full}-lambda-policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
         "Effect": "Allow",
         "Action": [
            "kms:Decrypt",
            "kms:GenerateDataKey",
            "s3:DeleteObject",
            "s3:ListBucket",
            "s3:HeadObject",
            "s3:GetObject",
            "s3:GetObjectVersion",
            "s3:PutObject",
            "s3:PutObjectAcl"
         ],
         "Resource": [
            "${data.aws_s3_bucket.cdn.arn}",
            "${data.aws_s3_bucket.cdn.arn}/*",
            "${data.aws_s3_bucket.storage.arn}",
            "${data.aws_s3_bucket.storage.arn}/*"
         ]
      },
      {
        "Effect": "Allow",
        "Action": [
          "logs:*"
        ],
        "Resource":  [
          "arn:aws:logs:*:*:*"
        ]
      }
  ]
}
EOF
}

resource "aws_iam_role_policy" "sm_policy" {
  name = "${local.env_ns_full}-sm-policy"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetSecretValue",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}



### API Gateway
resource "aws_iam_role" "api_gateway_account_role" {
  name = "${local.env_ns_full}-api-gateway-account-role"
  
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "",
        "Effect" : "Allow",
        "Principal" : {
          "Service": ["apigateway.amazonaws.com","lambda.amazonaws.com"]
        },
        "Action" : "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_policy" "api_gateway_cloudwatch_policy" {
  name = "${local.env_ns_full}-api-gateway-cloudwatch-policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams",
          "logs:PutLogEvents",
          "logs:GetLogEvents",
          "logs:FilterLogEvents"
        ],
        "Resource" : "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_gateway_s3" {
  role       = aws_iam_role.api_gateway_account_role.name
  policy_arn = aws_iam_policy.lambda.arn
}


resource "aws_api_gateway_account" "api_gateway_account" {
  cloudwatch_role_arn = aws_iam_role.api_gateway_account_role.arn
}



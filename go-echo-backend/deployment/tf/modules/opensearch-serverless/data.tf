locals {
  function_name = "${var.name}-${var.env}-central-logs"
  image_name    = "${var.name}-${var.env}-central-logs"
  image_tag     = "latest"
  repo_uri      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com"
}

data "aws_route53_zone" "opensearch" {
  name = var.hosted_zone_name
}

data "aws_acm_certificate" "opensearch" {
  domain = var.hosted_zone_name
}

data "aws_vpc" "selected" {
  id = var.vpc_id
}


data "aws_caller_identity" "current" {

}
data "aws_region" "current" {

}
data "aws_subnets" "private" {
  filter {
    name   = "tag:Name"
    values = ["*private*"]
  }
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.selected.id]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "tag:Name"
    values = ["*public*"]
  }
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.selected.id]
  }
}
data "aws_iam_policy_document" "access_policies" {
  statement {
    effect = "Allow"

    actions = [
      "aoss:*",
      "iam:ListUsers",
      "iam:ListRoles",
      "ec2:DescribeNetworkInterfaces",
      "ec2:CreateNetworkInterface",
      "ec2:DeleteNetworkInterface",
      "ec2:DescribeInstances",
      "ec2:AttachNetworkInterface"
    ]

    resources = ["*"]
  }
}

data "aws_iam_policy_document" "master_user_policy_document" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_log_and_invoke_policy" {

  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]

  }

  statement {
    effect = "Allow"

    actions = ["lambda:InvokeFunction"]

    resources = ["arn:aws:lambda:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:function:*"]
  }

}

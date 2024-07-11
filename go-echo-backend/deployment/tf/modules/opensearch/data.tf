data "aws_vpc" "selected" {
  id = var.vpc_id
}

data "aws_iam_roles" "role" {
  name_regex = "AWSServiceRoleForAmazonOpenSearchService"
}

locals {
  role_exists = can(data.aws_iam_roles.role.arns)
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


data "aws_iam_policy_document" "access_policies" {
  statement {
    actions = [
      "es:*"
    ]
    resources = [
      "arn:aws:es:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:domain/${local.namespace}/*"
    ]
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
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
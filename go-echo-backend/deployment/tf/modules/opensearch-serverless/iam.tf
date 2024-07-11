resource "aws_iam_role" "master_user_role" {
  name               = "${local.namespace}-opensearch-serverless"
  assume_role_policy = data.aws_iam_policy_document.master_user_policy_document.json
}

resource "aws_iam_role_policy" "frontend_lambda_role_policy" {
  name   = "${local.namespace}-lambda-role-policy"
  role   = "${aws_iam_role.master_user_role.id}"
  policy = "${data.aws_iam_policy_document.lambda_log_and_invoke_policy.json}"
}


resource "aws_iam_role_policy" "lambda_opensearch_role_policy" {
  name   = "${local.namespace}-lambda-opensearch-role-policy"
  role   = "${aws_iam_role.master_user_role.id}"
  policy = "${data.aws_iam_policy_document.access_policies.json}"
}


resource "aws_iam_role_policy_attachment" "route53" {
  role       = aws_iam_role.master_user_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonRoute53FullAccess"
}

resource "aws_iam_role_policy_attachment" "vpc" {
  role       = aws_iam_role.master_user_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
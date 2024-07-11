resource "aws_iam_role" "master_user_role" {
  name               = "${local.namespace}-master-user-role"
  assume_role_policy = data.aws_iam_policy_document.master_user_policy_document.json
}


resource "aws_iam_role_policy_attachment" "master_user_attachment" {
  role       = aws_iam_role.master_user_role.id
  policy_arn = "arn:aws:iam::aws:policy/AmazonOpenSearchServiceFullAccess"
}

resource "aws_iam_role_policy" "frontend_lambda_role_policy" {
  name   = "${local.namespace}-lambda-role-policy"
  role   = "${aws_iam_role.master_user_role.id}"
  policy = "${data.aws_iam_policy_document.lambda_log_and_invoke_policy.json}"
}

resource "aws_iam_role_policy_attachment" "AWSLambdaVPCAccessExecutionRole" {
    role       = aws_iam_role.master_user_role.id
    policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
resource "aws_iam_role" "codebuild" {
  name               = "${local.env_ns}-codebuild"
  assume_role_policy = "${data.aws_iam_policy_document.assume_by_codebuild.json}"
}
resource "aws_iam_role" "codedeploy" {
  name               = "${local.env_ns}-codedeploy"
  assume_role_policy = "${data.aws_iam_policy_document.assume_by_codedeploy.json}"
}

resource "aws_iam_role_policy" "codedeploy_policy" {
  name   = "${local.env_ns}-codedeploy_policy"
  role   = aws_iam_role.codedeploy.id
  policy = data.aws_iam_policy_document.codedeploy_policy.json
}

resource "aws_iam_role_policy" "codebuild" {
  role   = "${aws_iam_role.codebuild.name}"
  policy = "${data.aws_iam_policy_document.codebuild.json}"
}

resource "aws_iam_role" "pipeline" {
  name = "${local.env_ns}-pipeline-ecs-service-role"
  assume_role_policy = "${data.aws_iam_policy_document.assume_by_pipeline.json}"
}

resource "aws_iam_role_policy" "pipeline" {
  role = "${aws_iam_role.pipeline.name}"
  policy = "${data.aws_iam_policy_document.pipeline.json}"
}


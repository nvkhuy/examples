locals {
  env_ns_underscore = "${var.name}_${var.env}"
}

# This role has a trust relationship which allows
# to assume the role of ec2
resource "aws_iam_role" "ecs_exec" {
  name = "${local.env_ns_underscore}_ecs_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    },
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
  EOF
}

# This is a policy attachement for the "ecs" role, it provides access
# to the the ECS service.
resource "aws_iam_policy_attachment" "ecs_exec" {
  name = "${local.env_ns_underscore}_ecs_policy_document"
  roles = [aws_iam_role.ecs_exec.id]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_policy" "ecs_exec" {
  name = "${local.env_ns_underscore}_ecs_exec_policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeTags",
        "ecs:CreateCluster",
        "ecs:DeregisterContainerInstance",
        "ecs:DiscoverPollEndpoint",
        "ecs:Poll",
        "ecs:RegisterContainerInstance",
        "ecs:StartTelemetrySession",
        "ecs:UpdateContainerInstancesState",
        "ecs:Submit*",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "secretsmanager:GetSecretValue",
        "secretsmanager:PutSecretValue",
        "kms:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "ecs_exec" {
  name = "${local.env_ns_underscore}_ecs_instance_profile"
  role = aws_iam_role.ecs_exec.name
}
resource "aws_iam_role_policy_attachment" "ecs_role_attach" {
  role       = aws_iam_role.ecs_exec.name
  policy_arn = aws_iam_policy.ecs_exec.arn
}

/* S3 permission */

# [Data] IAM policy to define S3 permissions
# TF: https://www.terraform.io/docs/providers/aws/d/iam_policy_document.html
# AWS: http://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html
# AWS CLI: http://docs.aws.amazon.com/cli/latest/reference/iam/create-policy.html
data "aws_iam_policy_document" "s3_data_bucket_policy" {
  statement {
    sid = ""
    effect = "Allow"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:ListBucket",
      "s3:HeadBucket",
      "s3:GetBucketLocation"
    ]
    resources = [
      "arn:aws:s3:::${var.state_bucket}"
    ]
  }
  statement {
    sid = ""
    effect = "Allow"
    actions = [
      "s3:DeleteObject",
      "s3:GetObject",
      "s3:PutObject",
      "s3:PutObjectAcl"
    ]
    resources = [
      "arn:aws:s3:::${var.state_bucket}",
      "arn:aws:s3:::${var.state_bucket}/*",
      "arn:aws:s3:::${var.state_bucket}/${var.env}.env",
    ]
  }
}

# AWS IAM policy
# TF: https://www.terraform.io/docs/providers/aws/r/iam_policy.html
# AWS: http://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html
# AWS CLI: http://docs.aws.amazon.com/cli/latest/reference/iam/create-policy.html
resource "aws_iam_policy" "s3_policy" {
  name = "${local.env_ns_underscore}_s3_policy"
  policy = "${data.aws_iam_policy_document.s3_data_bucket_policy.json}"
}

# Attaches a managed IAM policy to an IAM role
# TF: https://www.terraform.io/docs/providers/aws/r/iam_role_policy_attachment.html
# AWS: http://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_managed-vs-inline.html
# AWS CLI: http://docs.aws.amazon.com/cli/latest/reference/iam/attach-role-policy.html
resource "aws_iam_role_policy_attachment" "ecs_role_s3_data_bucket_policy_attach" {
  role       = "${aws_iam_role.ecs_exec.name}"
  policy_arn = "${aws_iam_policy.s3_policy.arn}"
}


/*
ECS Task
*/

// AWS IAM Policy in format JSON for linked role ecs task
data "aws_iam_policy_document" "ecs_task" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"

      identifiers = [
        "ecs-tasks.amazonaws.com",
        "s3.amazonaws.com",
      ]
    }
  }
}

// - IAM role that the Amazon ECS container agent and the Docker daemon can assume
data "aws_iam_policy_document" "ecs_task_exec" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"

      identifiers = [
        "ecs-tasks.amazonaws.com",
        "s3.amazonaws.com",
      ]
    }
  }
}

// AWS IAM Role used for ecs task
resource "aws_iam_role" "ecs_task" {
  name = "${local.env_ns_underscore}_ecs_task"
  assume_role_policy = data.aws_iam_policy_document.ecs_task.json
}

resource "aws_iam_role_policy_attachment" "s3_task" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_role_policy_attachment" "sqs_task" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task" {
  role       = aws_iam_role.ecs_task.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonECS_FullAccess"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy_attachment" {
  role       = aws_iam_role.ecs_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "password_policy_secretsmanager" {
  name = "${local.env_ns_underscore}_ecs_task_secretsmanager"
  role = aws_iam_role.ecs_task.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": [
          "secretsmanager:GetSecretValue",
          "secretsmanager:PutSecretValue"
        ],
        "Effect": "Allow",
        "Resource": "*"
      }
    ]
  }
  EOF
}


resource "aws_iam_role_policy" "s3_storages" {
  name = "${local.env_ns_underscore}_s3_storages"
  role = aws_iam_role.ecs_task.id

  policy = <<-EOF
  {
    "Statement": [
        {
            "Action": [
                "s3:PutObjectAcl",
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject",
                "s3:ListObject",
                "s3:ListBucket"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:s3:::${var.datastore_bucket}",
                "arn:aws:s3:::${var.datastore_bucket}/*",
                "arn:aws:s3:::${var.storage_bucket}",
                "arn:aws:s3:::${var.storage_bucket}/*",
                "arn:aws:s3:::${var.cdn_bucket}",
                "arn:aws:s3:::${var.cdn_bucket}/*"
            ],
            "Sid": ""
        }
    ],
    "Version": "2012-10-17"
  }
  EOF
}

// AWS IAM Role used for ec2
resource "aws_iam_role" "ec2_s3_datastore" {
  name = "${local.env_ns_underscore}_ec2_s3_datastore"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "ec2_s3_datastore" {
  name = "${local.env_ns_underscore}_ec2_s3_datastore"
  role = aws_iam_role.ec2_s3_datastore.id

  policy = <<-EOF
  {
    "Statement": [
        {
            "Action": [
                "s3:PutObjectAcl",
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject",
                "s3:ListBucket"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:s3:::${var.datastore_bucket}",
                "arn:aws:s3:::${var.datastore_bucket}/*"
            ],
            "Sid": ""
        }
    ],
    "Version": "2012-10-17"
  }
  EOF
}

# Create an IAM instance profile
resource "aws_iam_instance_profile" "ec2_s3_profile" {
  name = "${local.env_ns_underscore}_ec2_s3_profile"
  role = aws_iam_role.ec2_s3_datastore.name
}

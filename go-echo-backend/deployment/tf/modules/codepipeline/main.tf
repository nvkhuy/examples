locals {
  env_ns = "${var.name}-${var.env}"

}
resource "aws_codebuild_project" "this" {
  for_each     = var.service_config
  name         = "${local.env_ns}-${each.value.name}-codebuild"
  description  = "Codebuild for the ECS Green/Blue ${local.env_ns} app"
  service_role = aws_iam_role.codebuild.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    image                       = lookup(each.value.codebuild_params, "image")
    type                        = lookup(each.value.codebuild_params, "type")
    compute_type                = lookup(each.value.codebuild_params, "compute_type")
    image_pull_credentials_type = lookup(each.value.codebuild_params, "cred_type")
    privileged_mode             = true

    environment_variable {
      name  = "SERVICE_NAME"
      value = each.value.name
    }

    environment_variable {
      name  = "BUILD_ENV"
      value = var.env
    }

    environment_variable {
      name  = "AWS_DEFAULT_REGION"
      value = var.region
    }

    environment_variable {
      name  = "AWS_PROFILE"
      value = var.profile
    }

  }

  source {
    type      = "CODEPIPELINE"
    buildspec = "buildspec.yml"
  }
}


resource "aws_codedeploy_app" "this" {
  for_each         = var.service_config
  compute_platform = "ECS"
  name             = "${local.env_ns}-${each.value.name}-service-deploy"
}

# resource "aws_codedeploy_deployment_group" "this" {
#   for_each               = var.service_config
#   app_name               = "${local.env_ns}-${each.value.name}-service-deploy"
#   deployment_group_name  = "${local.env_ns}-${each.value.name}-service-deploy-group"
#   deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"
#   service_role_arn       = aws_iam_role.codedeploy.arn

#   ecs_service {
#     cluster_name = var.cluster_name
#     service_name = "${local.env_ns}-${each.value.name}"
#   }

#   auto_rollback_configuration {
#     enabled = true
#     events  = ["DEPLOYMENT_FAILURE"]
#   }

#   deployment_style {
#     deployment_option = "WITH_TRAFFIC_CONTROL"
#     deployment_type   = "BLUE_GREEN"
#   }

#   blue_green_deployment_config {
#     deployment_ready_option {
#       action_on_timeout    = "CONTINUE_DEPLOYMENT"
#       wait_time_in_minutes = 0
#     }

#     terminate_blue_instances_on_deployment_success {
#       action                           = "TERMINATE"
#       termination_wait_time_in_minutes = 5
#     }
#   }

#   load_balancer_info {
#     target_group_pair_info {
#       prod_traffic_route {
#         listener_arns = var.https_listener_arns
#       }

#       test_traffic_route {
#         listener_arns = var.http_listener_arns
#       }
#       target_group {
#         name = var.target_groups_name_primary[each.value.name].arn
#       }

#       target_group {
#         name = var.target_groups_name_secondary[each.value.name].arn
#       }
#     }
#   }
# }

resource "aws_s3_bucket" "pipeline" {
  for_each = var.service_config
  bucket   = "${local.env_ns}-${each.value.name}-codepipeline-bucket"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "${local.env_ns}-${each.value.name}Codepipeline",
  "Statement": [
        {
            "Sid": "DenyUnEncryptedObjectUploads",
            "Effect": "Deny",
            "Principal": "*",
            "Action": "s3:PutObject",
            "Resource": "arn:aws:s3:::${local.env_ns}-${each.value.name}-codepipeline-bucket/*",
            "Condition": {
                "StringNotEquals": {
                    "s3:x-amz-server-side-encryption": "aws:kms"
                }
            }
        },
        {
            "Sid": "DenyInsecureConnections",
            "Effect": "Deny",
            "Principal": "*",
            "Action": "s3:*",
            "Resource": "arn:aws:s3:::${local.env_ns}-${each.value.name}-codepipeline-bucket/*",
            "Condition": {
                "Bool": {
                    "aws:SecureTransport": "false"
                }
            }
        }
    ]
}
POLICY
}

resource "aws_codepipeline" "this" {
  for_each = var.service_config
  name     = "${local.env_ns}-${each.value.name}-pipeline"
  role_arn = aws_iam_role.pipeline.arn

  artifact_store {
    location = "${local.env_ns}-${each.value.name}-codepipeline-bucket"
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name     = "Source"
      category = "Source"
      owner    = "AWS"
      provider = "CodeStarSourceConnection"
      version  = "1"

      output_artifacts = ["source_output"]

      configuration = {
        ConnectionArn        = var.codestar_arn
        FullRepositoryId     = lookup(each.value.codebuild_params, "git_repo")
        BranchName           = lookup(each.value.codebuild_params, "git_branch")
        OutputArtifactFormat = "CODE_ZIP"
      }
    }
  }

  stage {
    name = "Build"

    action {
      name     = "Build"
      category = "Build"
      owner    = "AWS"
      provider = "CodeBuild"
      version  = "1"

      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]

      configuration = {
        ProjectName = "${local.env_ns}-${each.value.name}-codebuild"
      }
    }
  }

  stage {
    name = "Deploy"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "ECS"
      input_artifacts = ["source_output"]
      version         = "1"

      configuration = {
        ClusterName = var.cluster_name
        ServiceName = "${local.env_ns}-${each.value.name}"
      }
    }
  }

  lifecycle {
    # prevent github OAuthToken from causing updates, since it's removed from state file
    ignore_changes = [stage[0].action[0].configuration]
  }

}

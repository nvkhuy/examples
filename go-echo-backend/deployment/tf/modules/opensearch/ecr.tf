data "aws_ecr_repository" "application_ecr_repo" {
  name     = local.image_name
}

resource "aws_ecr_lifecycle_policy" "application_ecr_repo_policy" {
  repository =  local.image_name
  policy     = <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Keep last 10 images",
            "selection": {
                "tagStatus": "tagged",
                "tagPrefixList": ["${data.aws_ecr_repository.application_ecr_repo.name}"],
                "countType": "imageCountMoreThan",
                "countNumber": 10
            },
            "action": {
                "type": "expire"
            }
        }
    ]
}
EOF
}


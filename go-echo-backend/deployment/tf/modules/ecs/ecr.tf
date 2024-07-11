data "aws_ecr_repository" "application_ecr_repo" {
  for_each = var.service_config
  name     = each.value.image_name
}

resource "aws_ecr_lifecycle_policy" "application_ecr_repo_policy" {
  for_each   = data.aws_ecr_repository.application_ecr_repo
  repository = each.value.name
  policy     = <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Keep last 10 images",
            "selection": {
                "tagStatus": "tagged",
                "tagPrefixList": ["${each.value.name}"],
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

depends_on = [ data.aws_ecr_repository.application_ecr_repo ]
}


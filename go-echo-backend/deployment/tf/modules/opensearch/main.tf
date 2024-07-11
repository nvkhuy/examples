locals {
  namespace     = "${var.name}-${var.env}-logs"
  subnet_ids    = slice(data.aws_subnets.private.ids, 0, var.instance_count)
  master_user   = "${var.name}-${var.env}-admin"

  function_name = "${var.name}-${var.env}-central-logs"
  image_name    = "${var.name}-${var.env}-central-logs"
  image_tag     = "latest"
  repo_uri      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com"
}

resource "aws_security_group" "opensearch_security_group" {
  count       = var.inside_vpc ? 1 : 0

  name        = "${local.namespace}-sg"
  vpc_id      = data.aws_vpc.selected.id
  description = "Allow inbound HTTP traffic"

  ingress {
    description = "HTTP from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"

    cidr_blocks = [
      data.aws_vpc.selected.cidr_block,
    ]
  }

    ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [
      data.aws_vpc.selected.cidr_block,
    ]
  }
}

resource "random_password" "password" {
  count       = var.internal_user_database_enabled && var.master_password == "" ? 1 : 0

  length      = 32
  special     = false
  min_lower   = 1
  min_numeric = 1
  min_special = 1
  min_upper   = 1
}

resource "aws_ssm_parameter" "opensearch_master_user" {
  count       = var.internal_user_database_enabled ? 1 : 0

  name        = "/service/${var.name}-${var.env}-opensearch/MASTER_USER"
  description = "opensearch_password for ${var.name}-${var.env} domain"
  type        = "SecureString"
  value       = "${var.master_user_name},${coalesce(var.master_password, try(random_password.password[0].result, ""))}"
}

resource "aws_iam_service_linked_role" "es" {
  count            = var.create_linked_role && !local.role_exists ? 1 : 0
  aws_service_name = var.aws_service_name_for_linked_role
}

resource "time_sleep" "role_dependency" {
  create_duration = "20s"

  triggers = {
    role_arn       = try(aws_iam_role.cognito_es_role[0].arn, null),
    linked_role_id = try(local.role_exists ? data.aws_iam_roles.role.id : aws_iam_service_linked_role.es[0].id, "11111")
  }
}

resource "aws_opensearch_domain" "opensearch" {
  domain_name    = local.namespace
  engine_version = "OpenSearch_${var.engine_version}"

  cluster_config {
    dedicated_master_count   = var.dedicated_master_count
    dedicated_master_type    = var.dedicated_master_type
    dedicated_master_enabled = var.dedicated_master_enabled
    instance_type            = var.instance_type
    instance_count           = var.instance_count
    zone_awareness_enabled   = var.zone_awareness_enabled
    zone_awareness_config {
      availability_zone_count = var.zone_awareness_enabled ? length(local.subnet_ids) : null
    }
  }

  advanced_security_options {
     enabled                        = var.advanced_security_options_enabled
    internal_user_database_enabled = var.internal_user_database_enabled
    master_user_options {
      master_user_arn      = var.master_user_arn == "" ? try(aws_iam_role.authenticated[0].arn, null) : var.master_user_arn
      master_user_name     = var.internal_user_database_enabled ? var.master_user_name : ""
      master_user_password = var.internal_user_database_enabled ? coalesce(var.master_password, try(random_password.password[0].result, "")) : ""
    }
  }

  dynamic "cognito_options" {
    for_each = var.cognito_enabled ? [1] : []
    content {
      enabled          = var.cognito_enabled
      user_pool_id     = aws_cognito_user_pool.user_pool[0].id
      identity_pool_id = aws_cognito_identity_pool.identity_pool[0].id
      role_arn         = time_sleep.role_dependency.triggers["role_arn"]
    }
  }

  software_update_options {
    auto_software_update_enabled = var.auto_software_update_enabled
  }
  domain_endpoint_options {
    enforce_https                   = var.custom_endpoint_enabled
    custom_endpoint_enabled         = var.custom_endpoint_enabled
    custom_endpoint                 = var.custom_endpoint_enabled ? var.custom_endpoint : null
    custom_endpoint_certificate_arn = var.custom_endpoint_enabled ? var.custom_endpoint_certificate_arn : null
    tls_security_policy             = var.tls_security_policy
  }

  ebs_options {
    ebs_enabled = var.ebs_enabled
    volume_size = var.ebs_volume_size
    volume_type = var.volume_type
    throughput  = var.throughput
  }


  node_to_node_encryption {
    enabled = var.node_to_node_encryption_enabled
  }

  encrypt_at_rest {
    enabled = var.encrypt_at_rest_enabled
  }
  dynamic "vpc_options" {
    for_each = var.inside_vpc ? [1] : []
    content {
      subnet_ids         = local.subnet_ids
      security_group_ids = [aws_security_group.es[0].id]
    }
  }

  dynamic "off_peak_window_options" {
    for_each = var.off_peak_window_start_time != null || var.off_peak_window_enabled != null ? [1] : []
    content {
      enabled = try(var.off_peak_window_enabled.enabled, null)
      dynamic "off_peak_window" {
        for_each = var.off_peak_window_start_time != null ? [1] : []
        content {
          window_start_time {
            hours   = var.off_peak_window_start_time.hours
            minutes = var.off_peak_window_start_time.minutes
          }
        }
      }
    }
  }

  access_policies = data.aws_iam_policy_document.access_policies.json

  depends_on = [aws_iam_service_linked_role.es[0], time_sleep.role_dependency]
}

resource "aws_route53_record" "opensearch_domain_record" {
  count = var.zone_id == "" ? 0 : 1
  zone_id = var.zone_id 
  name    = var.custom_endpoint
  type    = "CNAME"
  ttl     = "300"

  records = [aws_opensearch_domain.opensearch.endpoint]
  depends_on = [aws_opensearch_domain.opensearch]
}


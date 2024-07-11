provider "aws" {
  region  = var.region
  profile = var.profile
}

locals {
  public_alb_target_groups = { for service, config in var.service_config : service => config.alb_target_group }
  env_ns                   = "${var.name}-${var.env}"
}

module "iam" {
  source           = "../../modules/iam"
  env              = var.env
  name             = var.name
  state_bucket     = var.state_bucket
  datastore_bucket = var.datastore_bucket
  storage_bucket   = var.storage_bucket
  cdn_bucket       = var.cdn_bucket
}

module "vpc" {
  source             = "../../modules/vpc"
  name               = var.name
  env                = var.env
  cidr               = var.cidr
  availability_zones = var.availability_zones
  public_subnets     = var.public_subnets
  private_subnets    = var.private_subnets
}


module "rds_alb_security_group" {
  source      = "../../modules/security-group"
  name        = "${local.env_ns}-rds"
  description = "${local.env_ns}-rds"
  vpc_id      = module.vpc.vpc_id

  ingress_rules = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
  ]

  egress_rules = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]

}

module "ssh_security_group" {
  source      = "../../modules/security-group"
  name        = "${local.env_ns}-ssh"
  description = "${local.env_ns}-ssh"
  vpc_id      = module.vpc.vpc_id

  ingress_rules = [
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    },
  ]

  egress_rules = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
  ]

}

module "efs" {
  source  = "../../modules/efs"
  name    = var.name
  env     = var.env
  subnets = module.vpc.private_subnets
  vpc_id  = module.vpc.vpc_id
}


module "public_alb_security_group" {
  source      = "../../modules/security-group"
  name        = local.env_ns
  description = local.env_ns
  vpc_id      = module.vpc.vpc_id

  ingress_rules = var.public_alb_config.ingress_rules
  egress_rules  = var.public_alb_config.egress_rules
}


module "public_alb" {
  source            = "../../modules/alb"
  env               = var.env
  name              = local.env_ns
  subnets           = module.vpc.public_subnets
  vpc_id            = module.vpc.vpc_id
  target_groups     = local.public_alb_target_groups
  internal          = false
  listener_port     = 80
  listener_protocol = "HTTP"
  listeners         = var.public_alb_config.listeners
  security_groups   = [module.public_alb_security_group.security_group_id]
  certificate_arn   = var.certificate_arn
  hosted_zone_id    = data.aws_route53_zone.route53.id
}



module "ecs" {
  source                      = "../../modules/ecs"
  name                        = var.name
  env                         = var.env
  region                      = var.region
  service_config              = var.service_config
  ecs_task_execution_role_arn = module.iam.ecs_exec_role_arn
  ecs_task_role_arn           = module.iam.ecs_task_arn
  vpc_id                      = module.vpc.vpc_id
  private_subnets             = module.vpc.private_subnets
  public_alb_security_group   = module.public_alb_security_group
  public_alb_target_groups    = module.public_alb.target_groups
  account_id                  = data.aws_caller_identity.current.account_id
}

module "redis" {
  source = "../../modules/redis"

  env  = var.env
  name = var.name

  num_cache_clusters = 1

  private_subnets = module.vpc.private_subnets

  vpc_id = module.vpc.vpc_id

  node_type = "cache.t3.micro"
}
module "cdn" {
  source = "../../modules/cdn"

  env  = var.env
  name = var.name

  region                       = var.region
  profile                      = var.profile
  account_id                   = data.aws_caller_identity.current.account_id
  cdn_domain                   = var.cdn_domain
  cdn_s3_bucket                = var.cdn_s3_bucket
  certificate_arn              = var.certificate_arn_us_east_1
  storage_domain               = var.storage_domain
  storage_s3_bucket            = var.storage_s3_bucket
  access_control_allow_origins = var.access_control_allow_origins

  trending_domain    = var.trending_domain
  trending_s3_bucket = var.trending_s3_bucket

}

module "vpn" {
  source                  = "../../modules/vpn"
  env                     = var.env
  name                    = var.name
  alb_arn                 = module.public_alb.alb_arn
  target_groups           = var.vpn_config
  listeners               = var.vpn_listeners
  vpc_id                  = module.vpc.vpc_id
  instance_id             = var.vpn_instance_id
  certificate_arn         = var.certificate_arn
  whitelisted_cidr_blocks = ["${var.vpn_ip}/32", "${chomp(data.http.myip.body)}/32"]
  cidr_blocks             = [var.cidr]
}

resource "aws_route53_record" "vpn_record" {
  zone_id = data.aws_route53_zone.route53.id
  name    = var.vpn_domain
  type    = "A"
  alias {
    name                   = module.public_alb.alb_dns_name
    zone_id                = module.public_alb.alb_zone_id
    evaluate_target_health = true
  }
}


# module "opensearch_serverless" {
#   source           = "../../modules/opensearch-serverless"
#   vpc_id           = module.vpc.vpc_id
#   cidr_blocks      = [module.vpc.cidr_block]
#   name             = var.name
#   env              = var.env
#   collection_name  = "${var.name}-${var.env}-logs"
#   cloudwatch_logs  = var.cloudwatch_logs
#   profile          = var.profile
#   logs_domain      = var.logs_domain
#   hosted_zone_name = var.hosted_zone_name
#   users = ["Loi", "Huy", "Son"]

#   // https://aws.amazon.com/blogs/big-data/configure-saml-federation-for-amazon-opensearch-serverless-with-okta/
#   # saml_metadata_document_url = "https://dev-10782191.okta.com/app/exkfe4tb4gNWW14y05d7/sso/saml/metadata"
#   saml_provider_name         = "okta"
#   saml_group                 = "opensearch"
# }


module "opensearch" {
  source = "../../modules/opensearch"

  profile         = var.profile
  vpc_id          = module.vpc.vpc_id
  cidr_blocks     = [var.cidr]
  subnet_ids      = module.vpc.private_subnets
  name            = var.name
  env             = var.env
  engine_version  = "2.9"
  throughput      = 250
  ebs_volume_size = 20
  instance_type   = "t3.medium.search"
  instance_count  = 1

  cloudwatch_logs                   = var.cloudwatch_logs
  advanced_security_options_enabled = true
  cognito_enabled                   = true
  custom_endpoint                   = "logs.${var.hosted_zone_name}"
  zone_id                           = data.aws_route53_zone.route53.zone_id
  custom_endpoint_enabled           = true
  custom_endpoint_certificate_arn   = data.aws_acm_certificate.domain.arn
}

module "workshop-domain" {
  source                 = "../../modules/custom-domain"
  route53_id             = data.aws_route53_zone.route53.id
  alb_dns_name           = module.public_alb.alb_dns_name
  alb_zone_id            = module.public_alb.alb_zone_id
  alb_listener_http_arn  = module.public_alb.aws_alb_listener_http
  alb_listener_https_arn = module.public_alb.aws_alb_listener_https
  domain                 = "workshop.${var.hosted_zone_name}"
  redirect = {
    host = "lvlyrdpcbjb.typeform.com"
    path = "/to/fK0AuOsA"
  }
}
# resource "aws_route53_record" "root_record" {
#   zone_id = data.aws_route53_zone.route53.id
#   name    = var.hosted_zone_name
#   type    = "A"
#   alias {
#     name                   = module.public_alb.alb_dns_name
#     zone_id                = module.public_alb.alb_zone_id
#     evaluate_target_health = true
#   }
# }

resource "aws_alb_listener_rule" "alb_listener_rule" {
  listener_arn = module.public_alb.aws_alb_listener_http

  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
  condition {
    host_header {
      values = [var.hosted_zone_name]
    }

  }
}

resource "aws_alb_listener_rule" "alb_listener_rule_https" {
  listener_arn = module.public_alb.aws_alb_listener_https
  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = "www.${var.hosted_zone_name}"
    }
  }
  condition {
    host_header {
      values = [var.hosted_zone_name]
    }

  }
}

module "share-domain" {
  source                 = "../../modules/custom-domain"
  route53_id             = data.aws_route53_zone.route53.id
  alb_dns_name           = module.public_alb.alb_dns_name
  alb_zone_id            = module.public_alb.alb_zone_id
  alb_listener_http_arn  = module.public_alb.aws_alb_listener_http
  alb_listener_https_arn = module.public_alb.aws_alb_listener_https
  domain                 = "t.${var.hosted_zone_name}"
  redirect = {
    host = "api.joininflow.io"
    path = "/api/v1/common/share/link/#{path}"
  }
}

resource "aws_wafv2_web_acl_association" "WafWebAclAssociation" {
  resource_arn = module.public_alb.alb_arn
  web_acl_arn  = var.web_acl_arn
}

module "lambda_resize" {
  source = "../../modules/lambda"

  env               = var.env
  name              = var.name
  region            = var.region
  profile           = var.profile
  account_id        = data.aws_caller_identity.current.account_id
  cdn_domain        = var.cdn_domain
  cdn_s3_bucket     = var.cdn_s3_bucket
  storage_domain    = var.storage_domain
  storage_s3_bucket = var.storage_s3_bucket
  media_jwt_secret  = var.media_jwt_secret

  image_tag             = "latest"
  image_name            = "${var.name}-${var.env}-resize"
  function_name         = "resize"
  memory_size           = 3008
  concurrent_executions = null
  variables             = null
}

module "apigateway_resize" {
  source = "../../modules/apigateway"

  name          = var.name
  env           = var.env
  region        = var.region
  profile       = var.profile
  function_name = "resize"

  route = {
    path   = "/resize"
    method = "GET"
  }

  zone_id         = data.aws_route53_zone.route53.zone_id
  certificate_arn = data.aws_acm_certificate.domain.arn
  domain_name     = "resize.${var.hosted_zone_name}"
  depends_on      = [module.lambda_resize]
}


module "lambda_ffmpeg" {
  source = "../../modules/lambda"

  env               = var.env
  name              = var.name
  region            = var.region
  profile           = var.profile
  account_id        = data.aws_caller_identity.current.account_id
  cdn_domain        = var.cdn_domain
  cdn_s3_bucket     = var.cdn_s3_bucket
  storage_domain    = var.storage_domain
  storage_s3_bucket = var.storage_s3_bucket
  media_jwt_secret  = var.media_jwt_secret

  image_tag             = "latest"
  image_name            = "${var.name}-${var.env}-ffmpeg"
  function_name         = "ffmpeg"
  memory_size           = 3008
  concurrent_executions = null
  variables             = null
}

module "apigateway_ffmpeg" {
  source = "../../modules/apigateway"

  name          = var.name
  env           = var.env
  region        = var.region
  profile       = var.profile
  function_name = "ffmpeg"

  zone_id         = data.aws_route53_zone.route53.zone_id
  certificate_arn = data.aws_acm_certificate.domain.arn
  domain_name     = "ffmpeg.${var.hosted_zone_name}"
  route = {
    path   = "/ffmpeg"
    method = "GET"
  }

  depends_on = [module.lambda_ffmpeg]
}



module "lambda_rod" {
  source = "../../modules/lambda"

  env               = var.env
  name              = var.name
  region            = var.region
  profile           = var.profile
  account_id        = data.aws_caller_identity.current.account_id
  cdn_domain        = var.cdn_domain
  cdn_s3_bucket     = var.cdn_s3_bucket
  storage_domain    = var.storage_domain
  storage_s3_bucket = var.storage_s3_bucket
  media_jwt_secret  = var.media_jwt_secret

  image_tag             = "latest"
  image_name            = "${var.name}-${var.env}-rod"
  function_name         = "rod"
  memory_size           = 3008
  concurrent_executions = null
  variables             = null
}

module "apigateway_rod" {
  source = "../../modules/apigateway"

  name          = var.name
  env           = var.env
  region        = var.region
  profile       = var.profile
  function_name = "rod"

  zone_id         = data.aws_route53_zone.route53.zone_id
  certificate_arn = data.aws_acm_certificate.domain.arn
  domain_name     = "rod.${var.hosted_zone_name}"
  route = {
    path   = "/rod"
    method = "GET"
  }


  depends_on = [module.lambda_rod]
}

module "lambda_blur" {
  source = "../../modules/lambda"

  env               = var.env
  name              = var.name
  region            = var.region
  profile           = var.profile
  account_id        = data.aws_caller_identity.current.account_id
  cdn_domain        = var.cdn_domain
  cdn_s3_bucket     = var.cdn_s3_bucket
  storage_domain    = var.storage_domain
  storage_s3_bucket = var.storage_s3_bucket
  media_jwt_secret  = var.media_jwt_secret

  image_tag             = "latest"
  image_name            = "${var.name}-${var.env}-blur"
  function_name         = "blur"
  memory_size           = 1024
  concurrent_executions = null
  variables             = null
}

module "apigateway_blur" {
  source = "../../modules/apigateway"

  name          = var.name
  env           = var.env
  region        = var.region
  profile       = var.profile
  function_name = "blur"

  certificate_arn = data.aws_acm_certificate.domain.arn
  domain_name     = "blur.${var.hosted_zone_name}"
  zone_id         = data.aws_route53_zone.route53.zone_id
  route = {
    path   = "/blur"
    method = "GET"
  }

  depends_on = [module.lambda_blur]
}

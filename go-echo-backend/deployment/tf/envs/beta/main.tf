provider "aws" {
  region  = var.region
  profile = var.profile
}

locals {
  public_alb_target_groups = { for service, config in var.service_config : service => config.alb_target_group }
  env_ns                   = "${var.name}-${var.env}"
}

terraform {
  backend "s3" {}
}

data "terraform_remote_state" "state" {
  backend = "s3"
  config = {
    bucket         = "${var.state_bucket}"
    dynamodb_table = "${var.state_lock_table}"
    region         = "${var.region}"
    key            = "${var.env}.tfstate"
  }
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

module "public_alb_security_group" {
  source      = "../../modules/security-group"
  name        = local.env_ns
  description = local.env_ns
  vpc_id      = var.vpc_id

  ingress_rules = var.public_alb_config.ingress_rules
  egress_rules  = var.public_alb_config.egress_rules
}


module "public_alb" {
  source            = "../../modules/alb"
  env               = var.env
  name              = local.env_ns
  subnets           = data.aws_subnets.public.ids
  vpc_id            = var.vpc_id
  target_groups     = local.public_alb_target_groups
  internal          = false
  listener_port     = 80
  listener_protocol = "HTTP"
  listeners         = var.public_alb_config.listeners
  security_groups   = [module.public_alb_security_group.security_group_id]
  certificate_arn   = var.certificate_arn
  hosted_zone_id = data.aws_route53_zone.route53.zone_id
}

module "ecs" {
  source                      = "../../modules/ecs"
  name                        = var.name
  env                         = var.env
  region                      = var.region
  service_config              = var.service_config
  ecs_task_execution_role_arn = module.iam.ecs_exec_role_arn
  ecs_task_role_arn           = module.iam.ecs_task_arn
  vpc_id                      = var.vpc_id
  private_subnets             = data.aws_subnets.private.ids
  public_alb_security_group   = module.public_alb_security_group
  public_alb_target_groups    = module.public_alb.target_groups
  account_id                  = data.aws_caller_identity.current.account_id
}


module "share-domain" {
  source                 = "../../modules/custom-domain"
  route53_id             = data.aws_route53_zone.route53.zone_id
  alb_dns_name           = module.public_alb.alb_dns_name
  alb_zone_id            = module.public_alb.alb_zone_id
  alb_listener_http_arn  = module.public_alb.aws_alb_listener_http
  alb_listener_https_arn = module.public_alb.aws_alb_listener_https
  domain                 = "beta-t.${var.hosted_zone_name}"
  redirect = {
    host = "beta-api.joininflow.io"
    path = "/api/v1/common/share/link/#{path}"
  }
}

resource "aws_wafv2_web_acl_association" "WafWebAclAssociation" {
  resource_arn = module.public_alb.alb_arn
  web_acl_arn  = var.web_acl_arn
}


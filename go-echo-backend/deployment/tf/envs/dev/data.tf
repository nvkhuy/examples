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

data "aws_route53_zone" "route53" {
  name = var.hosted_zone_name
}

data "aws_acm_certificate" "domain" {
  domain = var.hosted_zone_name
}


data "aws_caller_identity" "current" {

}

data "http" "myip" {
  url = "http://ipv4.icanhazip.com"
}
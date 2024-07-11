data "aws_vpc" "selected" {
  id = var.vpc_id
}

data "aws_caller_identity" "current" {

}
data "aws_region" "current" {

}

data "aws_subnets" "private" {
  filter {
    name   = "tag:Name"
    values = ["*private*"]
  }
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.selected.id]
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "tag:Name"
    values = ["*public*"]
  }
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.selected.id]
  }
}

data "aws_route53_zone" "route53" {
  name = var.hosted_zone_name
}
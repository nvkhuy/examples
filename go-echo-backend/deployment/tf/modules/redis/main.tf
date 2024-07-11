locals {
  env_ns = "${var.name}-${var.env}"
}

// Docs: https://github.com/hashicorp/terraform-elasticache-example/blob/master/README.md
resource "aws_elasticache_subnet_group" "this" {
  name       = local.env_ns
  subnet_ids = var.private_subnets
}

resource "aws_security_group" "redis" {
  name   = "${local.env_ns}-redis"
  vpc_id = var.vpc_id

  ingress {
    protocol         = "tcp"
    from_port        = "6379"
    to_port          = "6379"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "${local.env_ns}"
  }
}


resource "aws_elasticache_replication_group" "this" {
  replication_group_id = local.env_ns
  description          = "${local.env_ns}-cluster"
  engine               = "redis"
  engine_version       = "7.1"
  node_type            = var.node_type
  port                 = 6379
  parameter_group_name = "default.redis7.cluster.on"

  snapshot_retention_limit = 5
  snapshot_window          = "00:00-05:00"

  subnet_group_name = aws_elasticache_subnet_group.this.name

  automatic_failover_enabled = true

  num_cache_clusters      = var.num_cache_clusters
  num_node_groups         = var.num_node_groups
  replicas_per_node_group = var.replicas_per_node_group

  security_group_ids = [aws_security_group.redis.id]

}                                                                  

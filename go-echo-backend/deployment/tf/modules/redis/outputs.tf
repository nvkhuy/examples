output "redis_arn" {
  value = aws_elasticache_replication_group.this.arn
}

output "redis_primary_endpoint_address" {
  value = aws_elasticache_replication_group.this.primary_endpoint_address
}

output "redis_configuration_endpoint_address" {
  value = aws_elasticache_replication_group.this.configuration_endpoint_address
}
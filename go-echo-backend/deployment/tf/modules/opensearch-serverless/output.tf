output "collection_enpdoint" {
  value = aws_opensearchserverless_collection.collection.collection_endpoint
}

output "dashboard_endpoint" {
  value = aws_opensearchserverless_collection.collection.dashboard_endpoint
}

output "saml" {
  value = aws_opensearchserverless_security_config.default
}

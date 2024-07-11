output "cloudfront_custom_domains" {
  value = aws_cloudfront_distribution.frontend.*.aliases
}

output "cloudfront_ids" {
  value = aws_cloudfront_distribution.frontend.*.id
}

output "cloudfront_domain" {
  value = aws_cloudfront_distribution.frontend.*.domain_name
}
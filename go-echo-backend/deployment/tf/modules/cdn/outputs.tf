output "storage_cloudfront_url" {
  value = aws_cloudfront_distribution.storage_distribution.aliases
}

output "cdn_cloudfront_url" {
  value = aws_cloudfront_distribution.cdn_distribution.aliases
}

output "storage_bucket" {
  value = data.aws_s3_bucket.storage.bucket
}

output "cdn_bucket" {
  value = data.aws_s3_bucket.cdn.bucket
}
output "cloudfront_cdn_distribution" {
  value = aws_cloudfront_distribution.cdn_distribution
}

output "cloudfront_storage_distribution" {
  value = aws_cloudfront_distribution.storage_distribution
}

output "cloudfront_storage_distribution_trusted_key_groups" {
  value = aws_cloudfront_distribution.storage_distribution.trusted_key_groups
}

resource "aws_s3_bucket" "frontend" {
  count  = length(var.webapp_config)
  bucket = var.webapp_config[count.index].bucket_name
}


resource "aws_s3_bucket_policy" "frontend" {
  count      = length(var.webapp_config)
  bucket     = aws_s3_bucket.frontend[count.index].id
  policy     = <<POLICY
{    
    "Version": "2012-10-17",    
    "Statement": [        
      {            
          "Sid": "PublicReadGetObject",            
          "Effect": "Allow",            
          "Principal": {
            "AWS": "*"
          },    
          "Action": [                
             "s3:GetObject"
          ],            
          "Resource": [
             "${aws_s3_bucket.frontend[count.index].arn}/*"        
          ]        
      }   
    ]
}
POLICY
  depends_on = [aws_s3_bucket.frontend]
}


resource "aws_s3_bucket_ownership_controls" "frontend" {
  count  = length(aws_s3_bucket.frontend)
  bucket = aws_s3_bucket.frontend[count.index].id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  count  = length(aws_s3_bucket.frontend)
  bucket = aws_s3_bucket.frontend[count.index].id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "frontend" {
  count  = length(aws_s3_bucket.frontend)
  bucket = aws_s3_bucket.frontend[count.index].id
  acl    = "public-read"

  depends_on = [
    aws_s3_bucket_ownership_controls.frontend,
    aws_s3_bucket_public_access_block.frontend,
  ]

}

resource "aws_s3_bucket_website_configuration" "frontend" {
  count  = length(aws_s3_bucket.frontend)
  bucket = aws_s3_bucket.frontend[count.index].id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}


resource "aws_cloudfront_distribution" "frontend" {
  count = length(aws_s3_bucket.frontend)
  origin {
    domain_name = aws_s3_bucket_website_configuration.frontend[count.index].website_endpoint
    origin_id   = aws_s3_bucket.frontend[count.index].bucket

    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
    }
  }
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  aliases = [var.webapp_config[count.index].domain]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_s3_bucket.frontend[count.index].bucket

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }

  viewer_certificate {
    acm_certificate_arn = var.certificate_arn
    ssl_support_method  = "sni-only"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  custom_error_response {
    error_caching_min_ttl = "0"
    error_code            = "403"
    response_code         = "200"
    response_page_path    = "/"
  }
  custom_error_response {
    error_caching_min_ttl = "0"
    error_code            = "404"
    response_code         = "200"
    response_page_path    = "/"
  }
  custom_error_response {
    error_caching_min_ttl = "0"
    error_code            = "400"
    response_code         = "200"
    response_page_path    = "/"
  }

}

resource "aws_route53_record" "frontend_record" {
  count   = length(var.webapp_config)
  zone_id = var.hosted_zone_id
  name    = var.webapp_config[count.index].domain
  type    = "A"
  alias {
    name                   = aws_cloudfront_distribution.frontend[count.index].domain_name
    zone_id                = aws_cloudfront_distribution.frontend[count.index].hosted_zone_id
    evaluate_target_health = false
  }
}


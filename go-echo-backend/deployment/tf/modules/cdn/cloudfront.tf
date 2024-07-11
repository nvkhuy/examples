resource "aws_cloudfront_origin_access_identity" "cdn" {
  comment = "Identity for S3 '${data.aws_s3_bucket.cdn.bucket}' bucket."
}

resource "aws_cloudfront_distribution" "cdn_distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  aliases             = [var.cdn_domain]

  origin {
    domain_name = data.aws_s3_bucket.cdn.bucket_regional_domain_name
    origin_id   = data.aws_s3_bucket.cdn.bucket

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.cdn.cloudfront_access_identity_path
    }
  }
  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = data.aws_s3_bucket.cdn.bucket
    viewer_protocol_policy = "redirect-to-https"
    compress               = true

    min_ttl     = 0
    default_ttl = 5 * 60
    max_ttl     = 60 * 60

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  viewer_certificate {
    acm_certificate_arn = var.certificate_arn
    ssl_support_method  = "sni-only"
  }


  custom_error_response {
    error_caching_min_ttl = 0
    error_code            = 404
    response_code         = 200
    response_page_path    = "/404.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

}

### For trending
resource "aws_cloudfront_origin_access_identity" "trending" {
  comment = "Identity for S3 '${data.aws_s3_bucket.trending.bucket}' bucket."
}

resource "aws_cloudfront_distribution" "trending_distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  aliases             = [var.trending_domain]

  origin {
    domain_name = data.aws_s3_bucket.trending.bucket_regional_domain_name
    origin_id   = data.aws_s3_bucket.trending.bucket

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.trending.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = data.aws_s3_bucket.trending.bucket
    viewer_protocol_policy = "redirect-to-https"

    compress = true

    min_ttl     = 0
    default_ttl = 5 * 60
    max_ttl     = 60 * 60

    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers_policy_storage.id


    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

  }

  viewer_certificate {
    acm_certificate_arn = var.certificate_arn
    ssl_support_method  = "sni-only"
  }

  custom_error_response {
    error_code            = "502"
    error_caching_min_ttl = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

}

### For storage
resource "aws_cloudfront_origin_access_identity" "storage" {
  comment = "Identity for S3 '${data.aws_s3_bucket.storage.bucket}' bucket."
}

resource "aws_cloudfront_distribution" "storage_distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  aliases             = [var.storage_domain]

  origin {
    domain_name = data.aws_s3_bucket.storage.bucket_regional_domain_name
    origin_id   = data.aws_s3_bucket.storage.bucket

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.storage.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = data.aws_s3_bucket.storage.bucket
    viewer_protocol_policy = "redirect-to-https"

    compress = true

    min_ttl     = 0
    default_ttl = 5 * 60
    max_ttl     = 60 * 60

    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers_policy_storage.id


    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

  }

  viewer_certificate {
    acm_certificate_arn = var.certificate_arn
    ssl_support_method  = "sni-only"
  }

  custom_error_response {
    error_code            = "502"
    error_caching_min_ttl = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

}

resource "aws_cloudfront_response_headers_policy" "security_headers_policy_storage" {
  name = "${var.name}-${var.env}-security-headers-policy"

  cors_config {
    access_control_allow_credentials = false
    access_control_allow_headers  {
      items = [
        "X-Custom-HTTP-Header"
      ]
    }
    access_control_allow_methods  {
      items = [
        "GET",
        "OPTIONS",
        "HEAD"
      ]
    }
    access_control_allow_origins  {
      items = var.access_control_allow_origins
    }
    access_control_expose_headers  {
      items = [
        "*"
      ]
    }
    access_control_max_age_sec = 600
    origin_override = true
  }
  
}
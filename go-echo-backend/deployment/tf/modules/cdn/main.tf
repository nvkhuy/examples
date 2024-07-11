provider "aws" {
  profile = var.profile
  region = var.region
}

locals {
  env_ns      = "${var.name}-${var.env}"
  env_ns_full = "${var.name}-${var.env}-resizer"
}


data "aws_s3_bucket" "storage" {
  bucket        = "${var.storage_s3_bucket}"
}

data "aws_s3_bucket" "cdn" {
  bucket        = "${var.cdn_s3_bucket}"
  # provisioner "local-exec" {
  #   command = "aws --region ${var.region} --profile ${var.profile} s3 cp test_image.png s3://${var.cdn_s3_bucket}/uploads --acl=public-read"
  # }
}

data "aws_s3_bucket" "trending" {
  bucket        = "${var.trending_s3_bucket}"
}

resource "aws_s3_bucket_public_access_block" "storage" {
  bucket = data.aws_s3_bucket.storage.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "trending" {
  bucket = data.aws_s3_bucket.trending.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "cdn" {
  bucket = data.aws_s3_bucket.cdn.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "cdn" {  
  bucket = data.aws_s3_bucket.cdn.id
policy = <<POLICY
{    
    "Version": "2012-10-17",    
    "Statement": [        
      {            
          "Sid": "PublicReadGetObject",            
          "Effect": "Allow",            
          "Principal": {
            "AWS": "${aws_cloudfront_origin_access_identity.cdn.iam_arn}"
          },    
          "Action": [                
             "s3:GetObject"            
          ],            
          "Resource": [
             "${data.aws_s3_bucket.cdn.arn}/*"            
          ]        
      }   
    ]
}
POLICY
}

resource "aws_s3_bucket_policy" "storage" {  
  bucket = data.aws_s3_bucket.storage.id
policy = <<POLICY
{    
    "Version": "2012-10-17",    
    "Statement": [        
      {            
          "Sid": "PublicReadGetObject",            
          "Effect": "Allow",            
          "Principal": {
            "AWS": "${aws_cloudfront_origin_access_identity.storage.iam_arn}"
          },            
          "Action": [                
             "s3:GetObject"            
          ],            
          "Resource": [
             "${data.aws_s3_bucket.storage.arn}/*"            
          ]        
      }   
    ]
}
POLICY
}


resource "aws_s3_bucket_policy" "trending" {  
  bucket = data.aws_s3_bucket.trending.id
policy = <<POLICY
{    
    "Version": "2012-10-17",    
    "Statement": [        
      {            
          "Sid": "PublicReadGetObject",            
          "Effect": "Allow",            
          "Principal": {
            "AWS": "${aws_cloudfront_origin_access_identity.trending.iam_arn}"
          },    
          "Action": [                
             "s3:GetObject"            
          ],            
          "Resource": [
             "${data.aws_s3_bucket.trending.arn}/*"            
          ]        
      }    
    ]
}
POLICY
}

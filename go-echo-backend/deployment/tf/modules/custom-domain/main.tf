
// custom
resource "aws_route53_record" "custom_record" {
  zone_id  = var.route53_id
  name     = var.domain
  type     = "A"
  alias {
    name                   = var.alb_dns_name
    zone_id                = var.alb_zone_id
    evaluate_target_health = true
  }
}

resource "aws_alb_listener_rule" "alb_listener_rule_custom_http" {
  listener_arn = var.alb_listener_http_arn
  
  action {
    type             = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  condition {
    dynamic "path_pattern" {
      for_each = var.path_pattern == null ? [] : [1]
      content {
        values = [var.path_pattern]
      }
    }
  }

  condition {
    host_header {
      values = [var.domain]
    }

  }
}

resource "aws_alb_listener_rule" "alb_listener_rule_custom_https" {
  listener_arn = var.alb_listener_https_arn
  action {
    type             = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host = var.redirect.host
      path = var.redirect.path
    }
  }



  
  condition {
    dynamic "path_pattern" {
      for_each = var.path_pattern == null ? [] : [1]
      content {
        values = [var.path_pattern]
      }
    }
  }

  condition {
    host_header {
      values = [var.domain]
    }
  }
}

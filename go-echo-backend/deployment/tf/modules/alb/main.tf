resource "aws_alb" "alb" {
  name               = var.name
  internal           = var.internal
  load_balancer_type = "application"
  subnets            = var.subnets
  security_groups    = var.security_groups
  enable_http2       = true
}

#Dynamically create the alb target groups for app services
resource "aws_alb_target_group" "alb_target_group" {
  for_each    = var.target_groups
  name        = "${lower(each.key)}-${var.env}"
  port        = each.value.port
  protocol    = each.value.protocol
  target_type = "ip"
  vpc_id      = var.vpc_id

  dynamic "health_check" {
    for_each = each.value.health_check == null ? [] : [true]

    content {
      matcher             = each.value.health_check.matcher
      path                = each.value.health_check.path
      interval            = each.value.health_check.interval
      timeout             = each.value.health_check.timeout
      unhealthy_threshold = each.value.health_check.unhealthy_threshold
    }
  }


  lifecycle {
    create_before_destroy = true
  }
}

# resource "aws_alb_target_group" "alb_target_group_secondary" {
#   for_each    = var.target_groups
#   name        = "${lower(each.key)}-${var.env}-secondary"
#   port        = each.value.port
#   protocol    = each.value.protocol
#   target_type = "ip"
#   vpc_id      = var.vpc_id

#   dynamic "health_check" {
#     for_each = each.value.health_check == null ? [] : [true]

#     content {
#       matcher             = each.value.health_check.matcher
#       path                = each.value.health_check.path
#       interval            = each.value.health_check.interval
#       timeout             = each.value.health_check.timeout
#       unhealthy_threshold = each.value.health_check.unhealthy_threshold
#     }
#   }


#   lifecycle {
#     create_before_destroy = true
#   }
# }


#Create the alb listener for the load balancer
resource "aws_alb_listener" "alb_listener" {
  for_each          = var.listeners
  load_balancer_arn = aws_alb.alb.id
  port              = each.value["listener_port"]
  protocol          = each.value["listener_protocol"]
  certificate_arn   = each.key == "HTTPS" ? var.certificate_arn : null
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "No routes defined"
      status_code  = "200"
    }
  }

}

#Create listener rules primary
resource "aws_alb_listener_rule" "alb_listener_rule" {
  for_each     = var.target_groups
  listener_arn = aws_alb_listener.alb_listener["HTTP"].arn
  # priority     = each.value.priority
  action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
  condition {
    host_header {
      values = each.value.host_header
    }

  }
}

resource "aws_alb_listener_rule" "alb_listener_rule_https" {
  for_each     = var.target_groups
  listener_arn = aws_alb_listener.alb_listener["HTTPS"].arn
  # priority     = each.value.priority
  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.alb_target_group[each.key].arn
  }
  condition {
    host_header {
      values = each.value.host_header
    }

  }
}


# resource "aws_alb_listener_rule" "alb_listener_rule_https_secondary" {
#   for_each     = var.target_groups
#   listener_arn = aws_alb_listener.alb_listener["HTTPS"].arn
#   # priority     = each.value.priority
#   action {
#     type             = "forward"
#     target_group_arn = aws_alb_target_group.alb_target_group_secondary[each.key].arn
#   }
#   condition {
#     host_header {
#       values = each.value.host_header
#     }

#   }
# }

resource "aws_route53_record" "frontend_record" {
  for_each = {
    for k, v in var.target_groups : k => v
    if v.domain != null
  }

  zone_id = var.hosted_zone_id
  name    = each.value.domain
  type    = "A"
  alias {
    name                   = aws_alb.alb.dns_name
    zone_id                = aws_alb.alb.zone_id
    evaluate_target_health = true
  }
}

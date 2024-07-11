locals {
  env_ns = "${var.name}-${var.env}"
}

data "aws_alb" "alb" {
  arn = var.alb_arn
}

#Dynamically create the alb target groups for app services
resource "aws_alb_target_group" "alb_target_group" {
  for_each    = var.target_groups
  name        = "${lower(each.key)}-${var.env}"
  port        = each.value.port
  protocol    = each.value.protocol
  target_type = "instance"
  vpc_id      = var.vpc_id

  dynamic "health_check" {
    for_each = each.value.health_check == null ? [] : [true]

    content {
      matcher             = each.value.matcher
      path                = each.value.path
      interval            = each.value.interval
      timeout             = each.value.timeout
      unhealthy_threshold = each.value.unhealthy_threshold
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}

#Create the alb listener for the load balancer
resource "aws_alb_listener" "alb_listener" {
  for_each          = var.listeners
  load_balancer_arn = data.aws_alb.alb.id
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

#Creat listener rules
resource "aws_alb_listener_rule" "alb_listener_rule" {
  for_each     = var.target_groups
  listener_arn = aws_alb_listener.alb_listener["HTTP"].arn
  
  action {
    type             = "redirect"
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

resource "aws_alb_target_group_attachment" "tgattachment" {
  for_each = var.target_groups

  target_group_arn = aws_alb_target_group.alb_target_group[each.key].arn
  target_id        = var.instance_id
}

resource "aws_security_group" "security_group" {
  name        = "${local.env_ns}-vpn-sg"
  description = "${local.env_ns}-vpn-sg"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = concat(var.whitelisted_cidr_blocks,var.cidr_blocks)
  }

  ingress {
    from_port   = 56789
    to_port     = 56789
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = var.cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

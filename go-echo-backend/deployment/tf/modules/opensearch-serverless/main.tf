locals {
  namespace              = "${var.name}-${var.env}"
  saml_metadata_document = "saml_metadata_document.xml"
  domain        = var.logs_domain
  custom_domain = "${local.domain}.${data.aws_route53_zone.opensearch.name}"

  users = [for user in var.users : "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${user}"]

  root_user   = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
  saml_groups = var.saml_group != "" ? ["saml/${data.aws_caller_identity.current.account_id}/${var.saml_provider_name}/group/${var.saml_group}"] : []
}

# Creates a security group
resource "aws_security_group" "security_group" {
  vpc_id = var.vpc_id
  tags = {
    Name = local.namespace
  }
}

resource "aws_security_group_rule" "allow_https" {
 type              = "ingress"
 description       = "HTTPS ingress"
 from_port         = 443
 to_port           = 443
 protocol          = "tcp"
 cidr_blocks       = var.cidr_blocks
 security_group_id = aws_security_group.security_group.id
}


# Allows all outbound traffic
resource "aws_vpc_security_group_egress_rule" "sg_egress" {
  security_group_id = aws_security_group.security_group.id

  cidr_ipv4   = "0.0.0.0/0"
  ip_protocol = "-1"
}

# Allows inbound traffic from within security group
resource "aws_vpc_security_group_ingress_rule" "sg_ingress" {
  security_group_id = aws_security_group.security_group.id

  referenced_security_group_id = aws_security_group.security_group.id
  ip_protocol                  = "-1"
}


# Creates an encryption security policy
resource "aws_opensearchserverless_security_policy" "encryption_policy" {
  name        = "${local.namespace}-encryption-policy"
  type        = "encryption"
  description = "encryption policy for ${local.namespace}"
  policy = jsonencode({
    Rules = [
      {
        Resource = [
          "collection/${var.collection_name}"
        ],
        ResourceType = "collection"
      }
    ],
    AWSOwnedKey = true
  })
}

# Creates a collection
resource "aws_opensearchserverless_collection" "collection" {
  name       = var.collection_name
  type       = "SEARCH"
  depends_on = [aws_opensearchserverless_security_policy.encryption_policy]
}

# Creates a network security policy
resource "aws_opensearchserverless_security_policy" "network_policy" {
  name        = "${local.namespace}-network-policy"
  type        = "network"
  description = "public access for dashboard, VPC access for collection endpoint"
  policy = jsonencode([
    {
      Description = "VPC access for collection endpoint",
      Rules = [
        {
          ResourceType = "collection",
          Resource = [
            "collection/${var.collection_name}"
          ]
        }
      ],
      AllowFromPublic = false,
      SourceVPCEs = [
        aws_opensearchserverless_vpc_endpoint.vpc_endpoint.id
      ]
    },
    {
      Description = "Public access for dashboards",
      Rules = [
        {
          ResourceType = "dashboard"
          Resource = [
            "collection/${var.collection_name}*"
          ]
        }
      ],
      AllowFromPublic = true
    }
  ])
}

# Creates a data access policy
resource "aws_opensearchserverless_access_policy" "data_access_policy" {
  name        = "${local.namespace}-data-access-policy"
  type        = "data"
  description = "allow index and collection access"
  policy = jsonencode([
    {
      Rules = [
        {
          ResourceType = "index",
          Resource = [
            "index/${var.collection_name}/*"
          ],
          Permission = [
            "aoss:*"
          ]
        },
        {
          ResourceType = "collection",
          Resource = [
            "collection/${var.collection_name}*"
          ],
          Permission = [
            "aoss:*"
          ]
        }
      ],
      Principal = concat([data.aws_caller_identity.current.arn, aws_iam_role.master_user_role.arn, local.root_user], local.users, local.saml_groups, )
    }
  ])
}
# Creates a VPC endpoint
resource "aws_opensearchserverless_vpc_endpoint" "vpc_endpoint" {
  name               = "${local.namespace}-vpc-endpoint"
  vpc_id             = var.vpc_id
  subnet_ids         = data.aws_subnets.private.ids
  security_group_ids = [aws_security_group.security_group.id]

}


resource "aws_opensearchserverless_lifecycle_policy" "default" {
  name = "${local.namespace}-retention"
  type = "retention"
  policy = jsonencode({
    "Rules" : [
      {
        "ResourceType" : "index",
        "Resource" : ["index/${var.collection_name}/*"],
        "MinIndexRetention" : "30d"
      }
    ]
  })
}

resource "aws_route53_record" "opensearch_domain_record" {
  zone_id = data.aws_route53_zone.opensearch.zone_id
  name    = local.custom_domain
  type    = "CNAME"
  ttl     = "300"

  records = [replace(aws_opensearchserverless_collection.collection.collection_endpoint,"https://","")]
}


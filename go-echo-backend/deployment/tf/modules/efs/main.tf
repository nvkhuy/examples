locals {
  env_ns = "${var.name}-${var.env}"
}


resource "aws_security_group" "nfs" {
  name   = "${local.env_ns}-nfs"
  vpc_id = var.vpc_id

  ingress {
    protocol         = "tcp"
    from_port        = "2049"
    to_port          = "2049"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    protocol         = "-1"
    from_port        = 0
    to_port          = 0
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "${local.env_ns}"
  }
}


# Creating Amazon EFS File system
resource "aws_efs_file_system" "myfilesystem" {
# Tagging the EFS File system with its value as Myfilesystem
  tags = {
    Name = "${local.env_ns}-efs"
  }

  encrypted = true
}

# Creating the EFS access point for AWS EFS File system
resource "aws_efs_access_point" "test" {
  file_system_id = aws_efs_file_system.myfilesystem.id
}

# Creating the AWS EFS System policy to transition files into and out of the file system.
resource "aws_efs_file_system_policy" "policy" {
  file_system_id = aws_efs_file_system.myfilesystem.id
# The EFS System Policy allows clients to mount, read and perform 
# write operations on File system 
# The communication of client and EFS is set using aws:secureTransport Option
  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Id": "Policy01",
    "Statement": [
        {
            "Sid": "Statement",
            "Effect": "Allow",
            "Principal": {
                "AWS": "*"
            },
            "Resource": "${aws_efs_file_system.myfilesystem.arn}",
            "Action": [
                "elasticfilesystem:ClientMount",
                "elasticfilesystem:ClientRootAccess",
                "elasticfilesystem:ClientWrite"
            ]
        }
    ]
}
POLICY
}
# Creating the AWS EFS Mount point in a specified Subnet 
# AWS EFS Mount point uses File system ID to launch.
resource "aws_efs_mount_target" "alpha" {
  count = length(var.subnets)
  file_system_id = aws_efs_file_system.myfilesystem.id
  subnet_id      = var.subnets[count.index]
  security_groups = [ aws_security_group.nfs.id ]
}
resource "aws_apprunner_service" "main" {
  service_name = "quotes-database"

  source_configuration {
    image_repository {
      image_configuration {
        port = "8080"

        # runtime_environment_secrets = {
        # TODO: DB Config
        # }
      }
      image_identifier      = var.image_url
      image_repository_type = "ECR"
    }
    auto_deployments_enabled = true

    # Access to ECR
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_ecr_access.arn
    }
  }

  health_check_configuration {
    path     = "/healthcheck"
    protocol = "HTTP"
  }

  instance_configuration {
    cpu    = "256"
    memory = "512"
  }

  network_configuration {
    ingress_configuration {
      is_publicly_accessible = true
    }
    egress_configuration {
      egress_type       = "VPC"
      vpc_connector_arn = aws_apprunner_vpc_connector.main.arn
    }
  }

  tags = {
    Name = "quotes-database"
  }
}

resource "aws_apprunner_vpc_connector" "main" {
  vpc_connector_name = "egress-to-vpc"
  subnets            = aws_subnet.private.*.id
  security_groups    = [aws_security_group.apprunner_egress.id]
}

resource "aws_security_group" "apprunner_egress" {
  name        = "apprunner-egress"
  description = "Security group for Apprunner VPC Connector"
  vpc_id      = aws_vpc.main.id
  tags = {
    Name = "apprunner-egress"
  }
}

resource "aws_vpc_security_group_egress_rule" "apprunner_to_database" {
  security_group_id = aws_security_group.apprunner_egress.id
  description       = "Apprunner VPC Connector to Database"
  from_port         = 5432
  to_port           = 5432
  ip_protocol       = "tcp"
  cidr_ipv4         = local.private_block
}

resource "aws_iam_role" "apprunner_ecr_access" {
  name               = "apprunner-ecr-access-role"
  assume_role_policy = data.aws_iam_policy_document.apprunner_assume_role.json
}

data "aws_iam_policy_document" "apprunner_assume_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["build.apprunner.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role_policy" "apprunner_ecr_access_policy" {
  name   = "apprunner-ecr-access-policy"
  role   = aws_iam_role.apprunner_ecr_access.id
  policy = data.aws_iam_policy_document.apprunner_ecr_access_policy.json
}

data "aws_iam_policy_document" "apprunner_ecr_access_policy" {
  statement {
    actions = [
      "ecr:GetAuthorizationToken",
    ]
    resources = ["*"]
  }
  statement {
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage"
    ]
    resources = [aws_ecr_repository.main.arn]
  }
}

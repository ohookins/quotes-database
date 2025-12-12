resource "aws_iam_role" "apprunner_instance" {
  name               = "apprunner-instance"
  assume_role_policy = data.aws_iam_policy_document.apprunner_instance_assume_role.json
}

data "aws_iam_policy_document" "apprunner_instance_assume_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["tasks.apprunner.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role_policy" "apprunner_instance_policy" {
  name   = "apprunner-instance-policy"
  role   = aws_iam_role.apprunner_instance.id
  policy = data.aws_iam_policy_document.apprunner_instance_policy.json
}

data "aws_iam_policy_document" "apprunner_instance_policy" {
  statement {
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.quotes_database_dsn.arn]
  }
}


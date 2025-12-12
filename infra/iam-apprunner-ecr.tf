resource "aws_iam_role" "apprunner_ecr_access" {
  name               = "apprunner-ecr-access-role"
  assume_role_policy = data.aws_iam_policy_document.apprunner_access_assume_role.json
}

data "aws_iam_policy_document" "apprunner_access_assume_role" {
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


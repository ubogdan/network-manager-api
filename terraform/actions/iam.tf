resource "aws_iam_user" "actions" {
  name = "github-actions-${var.app_name}-${var.stack_env}"
}

resource "aws_iam_access_key" "actions" {
  user = aws_iam_user.actions.name
  pgp_key = var.pgp_key
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_iam_user_policy" "actions_authorization" {
  name = "ecr-authorization-${var.app_name}-${var.stack_env}"
  user = aws_iam_user.actions.name
  policy = <<EOF
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Sid":"GetAuthorizationToken",
         "Effect":"Allow",
         "Action":[
            "ecr:GetAuthorizationToken"
         ],
         "Resource":"*"
      }
   ]
}
EOF
}

resource "aws_iam_user_policy" "actions-push" {
  name = "ecr-push-${var.app_name}-${var.stack_env}"
  user = aws_iam_user.actions.name

  policy = <<EOF
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Sid":"AllowPush",
         "Effect":"Allow",
         "Action":[
            "ecr:GetDownloadUrlForLayer",
            "ecr:BatchGetImage",
            "ecr:BatchCheckLayerAvailability",
            "ecr:PutImage",
            "ecr:InitiateLayerUpload",
            "ecr:UploadLayerPart",
            "ecr:CompleteLayerUpload"
         ],
         "Resource":"arn:aws:ecr:${var.region}:${data.aws_caller_identity.current.account_id}:repository/${var.app_name}"
      }
   ]
}
EOF
}

resource "aws_iam_user_policy" "actions-lambda" {
  name = "lambda-update-${var.app_name}-${var.stack_env}"
  user = aws_iam_user.actions.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:ListRoles",
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration"
      ],
      "Resource": "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:${var.app_name}-${var.stack_env}"
    }
  ]
}
EOF
}


output "aws_iam_api" {
  value = aws_iam_access_key.actions.id
}

output "aws_iam_secret" {
  value = aws_iam_access_key.actions.encrypted_secret
  sensitive = false
}
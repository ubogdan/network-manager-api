provider "aws" {
  region = var.stack_region
}

terraform {
  required_version = ">= 0.14"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }

  backend "s3" {
    encrypt = true
    region  = "eu-central-1"
  }
}

data "aws_caller_identity" "current" {}

resource "aws_lambda_function" "lambda_handler" {

  function_name = "${var.app_name}-${var.stack_env}"
  description   = "Lambda function for ${var.app_name}-${var.stack_env}"

  publish                        = true
  package_type                   = "Image"
  image_uri                      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.stack_region}.amazonaws.com/${var.app_name}:${var.app_version}"
  role                           = aws_iam_role.lambda.arn
  memory_size                    = var.memory
  timeout                        = var.timeout
  reserved_concurrent_executions = var.reserved_concurrency

  image_config {
    command = ["/app/service"]
  }

  environment {
    variables = {
    S3_BUCKET_REGION = "eu-central-1" }
  }

  tags = {
    Name        = "${var.app_name}-${var.stack_env}"
    Environment = var.stack_env
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${aws_lambda_function.lambda_handler.function_name}"

  retention_in_days = 7
  tags = {
    Environment = var.stack_env
    Service     = var.app_name
  }
}

resource "aws_iam_role_policy_attachment" "execution_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_sns_topic" "sns_topic" {
  name = "nm-email-notification"

  kms_master_key_id = "alias/aws/sns"

  tags = {
    Service     = var.app_name
    Environment = var.stack_env
  }
}

resource "aws_sqs_queue" "sns_topic" {
  name = "${aws_sns_topic.sns_topic.name}-dlq"

  message_retention_seconds = 604800

  kms_master_key_id = "alias/aws/sqs"

  tags = {
    Service     = var.app_name
    Environment = var.stack_env
  }
}

data "aws_iam_policy_document" "sns_topic_drl" {
  statement {
    actions = ["sqs:SendMessage"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      aws_sqs_queue.sns_topic.arn,
    ]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values = [
        aws_sns_topic.sns_topic.arn,
      ]
    }
  }
}

resource "aws_sqs_queue_policy" "sns_publish_policy" {
  queue_url = aws_sqs_queue.sns_topic.id
  policy    = data.aws_iam_policy_document.sns_topic_drl.json
}

resource "aws_sns_topic_subscription" "topic_subscription" {
  topic_arn      = aws_sns_topic.sns_topic.arn
  protocol       = "lambda"
  endpoint       = aws_lambda_function.lambda_handler.arn
  redrive_policy = jsonencode({ "deadLetterTargetArn" = aws_sqs_queue.sns_topic.arn })
}

data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${var.app_name}-${var.stack_env}"
  assume_role_policy = data.aws_iam_policy_document.lambda.json
}

data "aws_iam_policy_document" "ses_policy" {
  statement {
    actions = [
      "ses:SendEmail",
      "ses:SendRawEmail",
    ]

    resources = [
      "arn:aws:ses:${var.stack_region}:${data.aws_caller_identity.current.account_id}:identity/lfpanels.com",
    ]
  }
}
resource "aws_iam_policy" "ses_access" {
  name   = "ses-${var.app_name}-${var.stack_env}"
  policy = data.aws_iam_policy_document.ses_policy.json
}

resource "aws_iam_role_policy_attachment" "ses_access" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.ses_access.arn
}

resource "aws_lambda_permission" "sns_execute" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_handler.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.sns_topic.arn
}


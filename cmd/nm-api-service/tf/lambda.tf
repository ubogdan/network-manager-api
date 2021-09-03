data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_lambda_function" "lambda_handler" {

  function_name                  = "${var.app_name}-${var.stack_env}"
  description                    = "Lambda handler for ${aws_api_gateway_rest_api.gateway.name}"

  publish                        = true
  package_type                   = "Image"
  image_uri = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com/${var.app_name}:${var.app_version}"
  role                           = aws_iam_role.lambda.arn
  memory_size                    = var.memory
  timeout                        = var.timeout
  reserved_concurrent_executions = var.reserved_concurrency

  image_config {
    command = ["/app/service"]
  }

  environment {
    variables = {
      S3_BUCKET_REGION = "eu-central-1"
      AUTHORIZED_KEY = var.authorized_key
      API_BEARER_AUTH = var.bearer_auth
      API_BASE_PATH = var.base_path
      LICENSE_KEY = var.license_key
    }
  }

  tags = {
    Name           = "${var.app_name}-${var.stack_env}"
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name = "/aws/lambda/${aws_lambda_function.lambda_handler.function_name}"

  retention_in_days = 30
  tags = {
    Environment = var.stack_env
    Service = var.app_name
  }
}

resource "aws_api_gateway_base_path_mapping" "v1" {
  api_id      = aws_api_gateway_rest_api.gateway.id
  stage_name  = var.stack_env
  domain_name = var.domain_name
  base_path   = var.base_path
}
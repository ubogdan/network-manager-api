provider "aws" {
  region = "eu-central-1"
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
    key     = "${var.app_name}-${var.stack_env}"
  }
}

locals {
  headers = {
    "Access-Control-Allow-Headers"     = "'${join(",", var.allow_headers)}'"
    "Access-Control-Allow-Methods"     = "'${join(",", var.allow_methods)}'"
    "Access-Control-Allow-Origin"      = "'${var.allow_origin}'"
    "Access-Control-Max-Age"           = "'${var.allow_max_age}'"
    "Access-Control-Allow-Credentials" = var.allow_credentials ? "'true'" : ""
  }

  # Pick non-empty header values
  header_values = compact(values(local.headers))

  # Pick names that from non-empty header values
  header_names = matchkeys(keys(local.headers), values(local.headers), local.header_values)

  # Parameter names for method and integration responses
  parameter_names = formatlist("method.response.header.%s", local.header_names)

  # Map parameter list to "true" values
  true_list = split("|", replace(join("|", local.parameter_names), "/[^|]+/", "true"))

  # Integration response parameters
  integration_response_parameters = zipmap(local.parameter_names, local.header_values)

  # Method response parameters
  method_response_parameters = zipmap(local.parameter_names, local.true_list)
}


resource "aws_api_gateway_rest_api" "gateway" {
  name        = "nm-api-${var.stack_env}"
  description = "Regional Rest API Gateway"

  disable_execute_api_endpoint = true

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "deployment" {
  depends_on = [
    aws_api_gateway_integration.mock,
    aws_api_gateway_integration.lambda,
  ]

  rest_api_id = aws_api_gateway_rest_api.gateway.id
  stage_name  = var.stack_env

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.proxy.id,
      aws_api_gateway_method.lambda.id,
      aws_api_gateway_integration.lambda.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  parent_id   = aws_api_gateway_rest_api.gateway.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "mock" {
  authorization = "NONE"
  http_method   = "OPTIONS"
  resource_id   = aws_api_gateway_resource.proxy.id
  rest_api_id   = aws_api_gateway_rest_api.gateway.id

  request_parameters = {
    "method.request.querystring.authuser" = false
    "method.request.querystring.code"     = false
    "method.request.querystring.page"     = false
    "method.request.querystring.perPage"  = false
    "method.request.querystring.prompt"   = false
    "method.request.querystring.scope"    = false
    "method.request.querystring.sortBy"   = false
    "method.request.querystring.state"    = false
  }
}

resource "aws_api_gateway_integration" "mock" {
  http_method = aws_api_gateway_method.mock.http_method
  resource_id = aws_api_gateway_resource.proxy.id
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  type        = "MOCK"

  request_templates = {
    "application/json" = "{ \"statusCode\": 200 }"
  }
}

resource "aws_api_gateway_integration_response" "mock" {
  resource_id = aws_api_gateway_resource.proxy.id
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  http_method = aws_api_gateway_method.mock.http_method
  status_code = 200

  response_parameters = local.integration_response_parameters

  response_templates = {
    "application/json" = ""
  }

  depends_on = [
    aws_api_gateway_integration.mock,
    aws_api_gateway_method_response.mock,
  ]
}

resource "aws_api_gateway_method_response" "mock" {
  resource_id = aws_api_gateway_resource.proxy.id
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  http_method = aws_api_gateway_method.mock.http_method
  status_code = 200

  response_parameters = local.method_response_parameters

  response_models = {
    "application/json" = "Empty"
  }

  depends_on = [
    aws_api_gateway_method.mock,
  ]
}

resource "aws_api_gateway_method" "lambda" {
  rest_api_id   = aws_api_gateway_rest_api.gateway.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_method_settings" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  stage_name  = var.stack_env
  method_path = "*/*"

  settings {
    # Set throttling values
    throttling_burst_limit = 10
    throttling_rate_limit  = 100
  }
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  resource_id = aws_api_gateway_method.lambda.resource_id
  http_method = aws_api_gateway_method.lambda.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_handler.invoke_arn
}

resource "aws_api_gateway_method_response" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  resource_id = aws_api_gateway_method.lambda.resource_id
  http_method = aws_api_gateway_method.lambda.http_method
  status_code = "200"

  depends_on = [
    aws_api_gateway_method.lambda,
  ]
}

resource "aws_api_gateway_integration_response" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gateway.id
  resource_id = aws_api_gateway_method.lambda.resource_id
  http_method = aws_api_gateway_method.lambda.http_method
  status_code = aws_api_gateway_method_response.lambda.status_code

  depends_on = [
    aws_api_gateway_integration.lambda,
    aws_api_gateway_method_response.lambda,
  ]
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_lambda_function" "lambda_handler" {

  function_name = "${var.app_name}-${var.stack_env}"
  description   = "Lambda handler for ${aws_api_gateway_rest_api.gateway.name}"

  publish                        = true
  package_type                   = "Image"
  image_uri                      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com/${var.app_name}:${var.app_version}"
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
      AUTHORIZED_KEY   = var.authorized_key
      API_BEARER_AUTH  = var.bearer_auth
      API_BASE_PATH    = var.base_path
      LICENSE_KEY      = var.license_key
      BACKUP_BUCKET    = var.backup_bucket
    }
  }

  tags = {
    Name = "${var.app_name}-${var.stack_env}"
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

resource "aws_api_gateway_base_path_mapping" "v1" {
  api_id      = aws_api_gateway_rest_api.gateway.id
  stage_name  = var.stack_env
  domain_name = var.domain_name
  base_path   = var.base_path
}

resource "aws_dynamodb_table" "license-table" {
  name           = "nm-licenses"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "HardwareID"

  attribute {
    name = "HardwareID"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    APP = "network-manager"
  }
}

resource "aws_dynamodb_table" "release-table" {
  name         = "nm-releases"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Version"

  attribute {
    name = "Version"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  tags = {
    APP = "network-manager"
  }
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


resource "aws_iam_role_policy_attachment" "execution_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "tracing_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

data "aws_iam_policy_document" "dynamodb_table_access" {
  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:Query",
      "dynamodb:Scan",
    ]

    resources = [
      "arn:aws:dynamodb:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:table/nm-licenses",
      "arn:aws:dynamodb:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:table/nm-releases",
    ]
  }
}

resource "aws_iam_policy" "dynamodb_table_access" {
  name        = "AWSLambdaDynamoDBAccess-${var.app_name}-${var.stack_env}"
  description = "Allows access to the nm-licenses table"
  policy      = data.aws_iam_policy_document.dynamodb_table_access.json
}

resource "aws_iam_role_policy_attachment" "dynamodb_table_access" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.dynamodb_table_access.arn
}

data "aws_iam_policy_document" "backup_access" {
  statement {
    actions = [
      "s3:PutObject",
    ]

    resources = [
      "arn:aws:s3:::${var.backup_bucket}",
      "arn:aws:s3:::${var.backup_bucket}/*",
    ]
  }
}

resource "aws_iam_policy" "backup_access" {
  name        = "AWSLambdaBackupAccess-${var.app_name}-${var.stack_env}"
  description = "Allows access to the backup bucket"
  policy      = data.aws_iam_policy_document.backup_access.json
}

resource "aws_iam_role_policy_attachment" "backup_access" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.backup_access.arn
}

data "aws_iam_policy_document" "ses_access" {
  statement {
    actions = [
      "ses:SendEmail",
      "ses:SendRawEmail",
    ]

    resources = [
      "arn:aws:ses:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:identity/${var.ses_domain}"
    ]
  }
}


resource "aws_iam_policy" "ses_access" {
  name        = "AWSLambdaSESAccess-${var.app_name}-${var.stack_env}"
  description = "Allows access to ses"
  policy      = data.aws_iam_policy_document.dynamodb_table_access.json
}

resource "aws_iam_role_policy_attachment" "ses_access" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.ses_access.arn
}

resource "aws_lambda_permission" "lambda_gateway_exec" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.gateway.execution_arn}/*/*/*"
}
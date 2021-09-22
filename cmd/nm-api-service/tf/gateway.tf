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
}

locals {
  headers = {
    "Access-Control-Allow-Headers" = "'${join(",", var.allow_headers)}'"
    "Access-Control-Allow-Methods" = "'${join(",", var.allow_methods)}'"
    "Access-Control-Allow-Origin" = "'${var.allow_origin}'"
    "Access-Control-Max-Age" = "'${var.allow_max_age}'"
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
    "method.request.querystring.authuser"         = false
    "method.request.querystring.code"             = false
    "method.request.querystring.page"             = false
    "method.request.querystring.perPage"          = false
    "method.request.querystring.prompt"           = false
    "method.request.querystring.scope"            = false
    "method.request.querystring.sortBy"           = false
    "method.request.querystring.state"            = false
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
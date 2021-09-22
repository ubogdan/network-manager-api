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
  role  = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "tracing_role" {
  role = aws_iam_role.lambda.name
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

resource "aws_lambda_permission" "lambda_gateway_exec" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_handler.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.gateway.execution_arn}/*/*/*"
}
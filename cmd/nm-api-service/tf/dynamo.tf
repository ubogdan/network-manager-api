resource "aws_dynamodb_table" "license-table" {
  name         = "nm-licenses"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "HardwareID"

  attribute {
    name = "HardwareID"
    type = "S"
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

  tags = {
    APP = "network-manager"
  }
}

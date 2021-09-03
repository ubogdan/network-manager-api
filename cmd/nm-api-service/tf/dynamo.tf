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
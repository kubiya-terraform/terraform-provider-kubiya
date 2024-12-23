terraform {
  required_providers {
    kubiya = {
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_integration" "integration" {
  name        = "mevrat-aws"
  description = "main aws account"

  configs = [
    {
      name       = "india"
      is_default = true
      vendor_specific = {
        arn    = "arn:aws:iam::590184027143:role/forkubiya"
        region = "ap-south-1"
      }
    },
    {
      name       = "brazil"
      is_default = false
      vendor_specific = {
        arn    = "arn:aws:iam::637423537751:role/brole1"
        region = "sa-east-1"
      }
    }
  ]
}

output "integration" {
  value = kubiya_integration.integration
}
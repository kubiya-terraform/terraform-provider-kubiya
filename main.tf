terraform {
  required_providers {
    kubiya = {
#       source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_source" "item" {
  url = "https://github.com/finebee/terraform-golden-usecases"
  dynamic_config = {
    michael = "hello"
  }
}

output "output" {
  value = kubiya_source.item
}

terraform {
  required_providers {
    kubiya = {
      # source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_source" "item" {
  dynamic_config = var.config_json
  url            = "https://github.com/kubiyabot/community-tools/tree/main/just_in_time_access"
}

output "output" {
  value = kubiya_source.item
}
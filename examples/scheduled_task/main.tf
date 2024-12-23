terraform {
  required_providers {
    kubiya = {
      #       source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_webhook" "mevrat" {
  name        = "mevrat"
  agent       = "Terraform-IaC"
  source      = "mevrat"
  prompt      = "mevrat"
  filter      = ""
  destination = "mevrat.avraham@kubiya.ai"
}

# resource "kubiya_scheduled_task" "item" {
#   repeat         = "daily"
#   channel_id     = "C08041WFAKT"
#   agent          = "mevrat_agent"
#   scheduled_time = "2024-12-01T05:00:00"
#   description    = "mevrat mevrat test task description"
# }

output "output" {
  value = kubiya_webhook.mevrat
}
terraform {
  required_providers {
    kubiya = {
      #       source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_scheduled_task" "example" {
  channel_id  = "C08041WFAKT"
  agent       = "mevrat_agent"
  repeat      = "*/15 * * * *"
  description = "mevrat mevrat test task description"
}

output "output" {
  value = kubiya_scheduled_task.example
}
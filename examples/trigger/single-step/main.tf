terraform {
  required_providers {
    kubiya = {
      # source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_trigger" "workflow_trigger" {
  name   = "workflow-trigger"
  runner = "core-testing-1"

  workflow = jsonencode({
    name    = "Echo via Webhook"
    version = 1
    steps = [
      {
        name = "echo"
        executor = {
          type = "command"
          config = {
            command = "echo \"Hello from webhook\""
          }
        }
      }
    ]
  })
}

output "trigger_url" {
  value       = kubiya_trigger.workflow_trigger.url
  description = "The webhook URL to trigger the workflow"
}

output "workflow_status" {
  value       = kubiya_trigger.workflow_trigger.status
  description = "The current status of the workflow"
}

output "workflow_id" {
  value       = kubiya_trigger.workflow_trigger.workflow_id
  description = "The ID of the created workflow"
}
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_trigger" "multi_step_workflow" {
  name   = "multi-step-workflow"
  runner = "kubiya-hosted"

  workflow = jsonencode({
    name    = "Multi-Step Workflow"
    version = 1
    steps = [
      {
        name = "step1"
        executor = {
          type = "command"
          config = {
            command = "echo \"Starting workflow...\""
          }
        }
      },
      {
        name = "step2"
        executor = {
          type = "command"
          config = {
            command = "date"
          }
        }
      },
      {
        name = "step3"
        executor = {
          type = "command"
          config = {
            command = "echo \"Workflow completed!\""
          }
        }
      }
    ]
  })
}

output "webhook_url" {
  value       = kubiya_trigger.multi_step_workflow.url
  description = "POST to this URL to trigger the workflow"
  sensitive   = true
}

output "workflow_info" {
  value = {
    id     = kubiya_trigger.multi_step_workflow.workflow_id
    status = kubiya_trigger.multi_step_workflow.status
    name   = kubiya_trigger.multi_step_workflow.name
  }
  description = "Information about the created workflow"
}
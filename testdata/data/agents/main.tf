terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

resource "kubiya_agent" "test_agent" {
  name        = var.agent_name
  runner      = var.runner_name
  description = var.description
  instructions = var.instructions

  # Optional fields with defaults
  model = var.model  # Defaults to "gpt-4o"
  is_debug_mode = var.debug_mode # Defaults to false

  # Access control
  users = var.users
  groups = var.groups

  # Integrations and capabilities
  sources = var.sources
  secrets = var.secrets
  integrations = var.integrations

  # Conversation starters
  starters = var.starters

  # Environment variables
  environment_variables = var.environment_variables
}

output "agent" {
  value = kubiya_agent.test_agent.id
}

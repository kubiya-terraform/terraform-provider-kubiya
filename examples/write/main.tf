terraform {
  required_providers {
    kubiya = {
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {
  user_key = "4534tgq3gq354gq5"
}

resource "kubiya_agent" "my_agent" {
  name            = "my_agent-name"
  image           = "my_agent-image"
  llm_model       = "my_agent-llm-model"
  description     = "my_agent-description"
  ai_instructions = "my_agent-ai-instructions"

  links                 = ["my_agent-link-1", "my_agent-link-2"]
  runners               = ["my_agent-runner-1", "my_agent-runner-2"]
  secrets               = ["my_agent-secret-1", "my_agent-secret-2"]
  starters              = ["my_agent-starter-1", "my_agent-starter-2"]
  integrations          = ["my_agent-integration-1", "my_agent-integration-2"]
  environment_variables = { "agent-2-env-1" = "my_agent-env-value-1", "my_agent-env-2" = "my_agent-env-value-2" }
}

output "resource_id" {
  value = kubiya_agent.my_agent
}
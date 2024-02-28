terraform {
  required_providers {
    kubiya = {
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {
  email        = ""
  user_key     = ""
  organization = ""
}

resource "kubiya_agent" "agent" {
  secrets = [
    "JFROG_ACCESS_TOKEN",
    "AWS_SECRET_ACCESS_KEY",
    "AWS_ACCESS_KEY_ID",
    "CONFIGCAT_TOKEN",
    "AWS_SESSION_TOKEN",
    "AWS_DEFAULT_REGION",
    "PULUMI_ACCESS_TOKEN"
  ]
  integrations = [
    "github",
    "jira",
    "kubernetes",
    "slack",
    "aws"
  ]
  links                 = [""]
  starters              = [""]
  environment_variables = {
    DEBUG     = "1"
    LOG_LEVEL = "INFO"
  }
  llm_model       = "azure/gpt-4"
  description     = "description"
  name            = "mevrat_agent_2"
  ai_instructions = "ai_instructions"
  runners         = ["aks-dev-tunnel"]
  image           = "kubiya/base-agent:latest"
}

output "default_agent_id" {
  value = kubiya_agent.agent
}

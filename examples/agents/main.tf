terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      #       source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {
  //Your Kubiya API Key will be taken from the
  //environment variable KUBIYA_API_KEY
  //To set the key, please use export KUBIYA_API_KEY="YOUR_API_KEY"
}

resource "kubiya_agent" "agent" {
  //Mandatory Fields
  name = "Mevrat-Camel-Case"             //String
  runner = "core-testing-1"           //String
  description  = "This agent can perform various Kubernetes tasks using kubectl"
  instructions = "You are an AI agent specialized in managing Kubernetes clusters using kubectl. Your tasks include monitoring pod statuses, scaling deployments, updating container images, and reporting issues to JIRA."

  //Optional fields, String
  #   model = "azure/gpt-4"  // If not provided, Defaults to "azure/gpt-4"
  //If not provided, Defaults to "ghcr.io/kubiyabot/kubiya-agent:stable"
  #   image = "ghcr.io/kubiyabot/kubiya-agent:stable"

  //Optional Fields (omitting will retain the current values):
  //Arrays
  integrations = ["github_app"]
  users = ["mevrat.avraham@kubiya.ai"]
  groups = ["Admin"]
}

output "agent" {
  value = kubiya_agent.agent
}

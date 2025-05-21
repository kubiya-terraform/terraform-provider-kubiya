terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {

}

resource "kubiya_agent" "agent" {

  //Mandatory Fields
  name = "Integration-Test-Agent"             //String
  runner = "runnerv2-5-vcluster"           //String
  description  = "This agent is a part of integration testing"
  instructions = "You are an AI agent specialized in managing Kubernetes clusters using kubectl. Your tasks include monitoring pod statuses, scaling deployments, updating container images, and reporting issues to JIRA."

  //Optional Fields (omitting will retain the current values):
  //Arrays
  integrations = ["github_app"]
  groups = ["Admin"]
}

output "agent" {
  value = kubiya_agent.agent
}

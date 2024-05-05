terraform {
  required_providers {
    kubiya = {
      source = "hashicorp.com/edu/Kubiya"
    }
  }
}

provider "kubiya" {
  user_key = "***REMOVED***"
}

resource "kubiya_agent" "bla_agent" {
  name = "bla"
}

output "resource_id" {
  value = kubiya_agent.bla_agent
}
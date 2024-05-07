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

data "kubiya_agent" "agent_1" {
  id = "agent-1-id"
}


output "default_agent_id" {
  value = data.kubiya_agent.agent_1.id
}
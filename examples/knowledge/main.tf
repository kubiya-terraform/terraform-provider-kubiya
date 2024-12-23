terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      # source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

resource "kubiya_knowledge" "knowledge" {
  // Required
  name    = "terraform-name-prod"
  groups = ["Admin"]
  content = "terraform-content-update"
  description = "terraform-description-prod"

  // Optional
  labels = ["label-1"]
  supported_agents = ["mevrat-enforcer"]
}

output "knowledge" {
  value = kubiya_knowledge.knowledge
}

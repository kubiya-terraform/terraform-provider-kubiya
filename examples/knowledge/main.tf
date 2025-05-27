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
}

output "knowledge" {
  value = kubiya_knowledge.knowledge
}

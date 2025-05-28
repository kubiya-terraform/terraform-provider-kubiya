terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
      # source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {}

# Complete example with all attributes
resource "kubiya_knowledge" "compliance_guidelines" {
  name        = "compliance-guidelines"
  description = "Company compliance and security guidelines"
  content     = "All cloud resources must be encrypted at rest and in transit. Access must be granted based on least privilege principles."

  groups = ["Admin"]
  labels = ["compliance", "security", "guidelines"]
  supported_agents = ["demo_teammate"]
}

output "knowledge" {
  value = kubiya_knowledge.compliance_guidelines
}

terraform {
  required_providers {
    kubiya = {
      # source = "kubiya-terraform/kubiya"
      source = "hashicorp.com/edu/kubiya"
    }
  }
}

provider "kubiya" {
  // API key is read from KUBIYA_API_KEY environment variable
}

# Example 1: Basic Slack integration with a single channel
resource "kubiya_external_knowledge" "example" {
  vendor = "slack"
  config = {
    channel_ids = ["C0735KZ7Z0A"]  # Slack channel ID
  }
}

# Output the external knowledge details
output "example_id" {
  value = kubiya_external_knowledge.example.id
}

output "example_details" {
  value = {
    id         = kubiya_external_knowledge.example.id
    vendor     = kubiya_external_knowledge.example.vendor
    org        = kubiya_external_knowledge.example.org
    start_date = kubiya_external_knowledge.example.start_date
    config     = kubiya_external_knowledge.example.config
  }
} 
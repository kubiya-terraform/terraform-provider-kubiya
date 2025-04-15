---
page_title: "kubiya_knowledge Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Knowledge resource to manage specific information for agents.
---

# kubiya_knowledge (Resource)

Provides a Kubiya Knowledge resource. This allows knowledge resources to be created, updated, and deleted on the Kubiya platform. Knowledge resources add specific information that can be used by Kubiya agents.

## Example Usage

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # Your Kubiya API Key will be taken from the 
  # environment variable KUBIYA_API_KEY
  # To set the key, please use export KUBIYA_API_KEY="YOUR_API_KEY"
}

# Simple knowledge resource
resource "kubiya_knowledge" "deployment_procedures" {
  name        = "deployment-procedures"
  description = "Standard procedures for application deployment"
  content     = file("${path.module}/deployment-procedures.md")
  format      = "markdown"
  
  # Groups are required
  groups      = ["DevOps"]
}

# Knowledge resource with access controls and labels
resource "kubiya_knowledge" "architecture_diagram" {
  name        = "architecture-diagram"
  description = "System architecture diagram and explanation"
  content     = file("${path.module}/architecture.md")
  format      = "markdown"
  
  # Specify which user groups can access this knowledge
  groups = ["Engineering", "DevOps", "Product"]
  
  # Add labels for organization and search
  labels = ["architecture", "documentation", "system-design"]
  
  # Specify which agents can use this knowledge
  supported_agents = ["infra-agent", "devops-assistant"]
}

# JSON format knowledge resource
resource "kubiya_knowledge" "api_schemas" {
  name        = "api-schemas"
  description = "API schemas for our services"
  content     = file("${path.module}/api-schemas.json")
  format      = "json"
  
  # Restrict access to admin groups only
  groups = ["Admin"]
  
  # Tag with relevant labels
  labels = ["api", "schema", "reference"]
}

# Complete example with all attributes
resource "kubiya_knowledge" "compliance_guidelines" {
  name        = "compliance-guidelines"
  description = "Company compliance and security guidelines"
  content     = "All cloud resources must be encrypted at rest and in transit. Access must be granted based on least privilege principles."
  format      = "text"
  
  groups = ["Admin", "Security", "Compliance"]
  labels = ["compliance", "security", "guidelines"]
  supported_agents = ["security-agent", "compliance-bot", "auditor"]
}

output "knowledge_id" {
  value = kubiya_knowledge.compliance_guidelines.id
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the knowledge resource.
* `description` - (Required) A description of the knowledge resource.
* `content` - (Required) The content of the knowledge resource. This can be provided directly or loaded from a file using the `file()` function.
* `groups` - (Required) List of user groups that have access to this knowledge resource.

### Optional Arguments

* `format` - (Optional) The format of the content. Values: "markdown", "text", "json", etc. Defaults to "text".
* `labels` - (Optional) List of labels/tags to help organize and find knowledge resources.
* `supported_agents` - (Optional) List of agent names that can access and use this knowledge resource. If not specified, all agents can use the knowledge.
* `type` - (Optional) The type of the knowledge. Defaults to "knowledge".

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the knowledge resource.
* `owner` - The owner of the knowledge resource.

## Access Control

Knowledge resources in Kubiya support fine-grained access control through:

1. **Groups**: Restrict access to specific user groups in your organization (required)
2. **Supported Agents**: Limit which AI Teammates (agents) can utilize this knowledge
3. **Labels**: Organize knowledge and improve discoverability

These controls help you manage who can see sensitive information and which agents can leverage specific knowledge in their responses.

## Import

Knowledge resources can be imported using the `id`:

```
$ terraform import kubiya_knowledge.example KNOWLEDGE_ID
``` 
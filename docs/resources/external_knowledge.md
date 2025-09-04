---
page_title: "kubiya_external_knowledge Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_external_knowledge resource manages external knowledge integrations in the Kubiya platform.
---

# kubiya_external_knowledge (Resource)

The `kubiya_external_knowledge` resource manages external knowledge integrations in the Kubiya platform. This resource allows you to connect external systems to the Kubiya RAG (Retrieval-Augmented Generation) system through a vendor-agnostic interface, enabling agents to access and learn from external data sources.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. Access to the external system you want to integrate (e.g., Slack workspace)
4. Appropriate permissions in the external system

## Example Usage

### 1. Basic Slack Integration

Connect a single Slack channel:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key is automatically read from KUBIYA_API_KEY environment variable
}

resource "kubiya_external_knowledge" "slack_general" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890"]  # General channel ID
  }
}

resource "kubiya_agent" "slack_aware_agent" {
  name         = "slack-aware-assistant"
  runner       = "kubiya-hosted"
  description  = "Agent with access to Slack knowledge"
  instructions = "You have access to historical Slack conversations. Reference relevant discussions when answering questions."
}
```

**Expected Outcome**: Creates an integration that allows agents to access knowledge from the specified Slack channel.

### 2. Multiple Slack Channels

Connect multiple Slack channels for broader knowledge access:

```hcl
resource "kubiya_external_knowledge" "slack_multi" {
  vendor = "slack"
  config = {
    channel_ids = [
      "C1234567890",  # #general
      "C0987654321",  # #engineering
      "C1111111111",  # #product
      "C2222222222"   # #support
    ]
  }
}

resource "kubiya_agent" "knowledge_agent" {
  name         = "knowledge-assistant"
  runner       = "kubiya-hosted"
  description  = "Agent with broad organizational knowledge"
  instructions = <<-EOT
    You have access to conversations from multiple Slack channels.
    Use this knowledge to:
    - Answer questions about past discussions
    - Provide context on decisions made
    - Reference relevant team communications
    - Help find information shared in Slack
  EOT
}
```

**Expected Outcome**: Creates an integration providing access to multiple Slack channels for comprehensive organizational knowledge.

### 3. Department-Specific Knowledge

Create department-specific knowledge integrations:

```hcl
resource "kubiya_external_knowledge" "engineering_slack" {
  vendor = "slack"
  config = {
    channel_ids = [
      "C3333333333",  # #eng-general
      "C4444444444",  # #eng-backend
      "C5555555555"   # #eng-frontend
    ]
  }
}

resource "kubiya_external_knowledge" "product_slack" {
  vendor = "slack"
  config = {
    channel_ids = [
      "C6666666666",  # #product-general
      "C7777777777",  # #product-design
      "C8888888888"   # #product-analytics
    ]
  }
}

resource "kubiya_agent" "engineering_agent" {
  name         = "engineering-knowledge-agent"
  runner       = "kubiya-hosted"
  description  = "Engineering team knowledge assistant"
  instructions = "You have access to engineering team Slack channels. Help with technical discussions and decisions."
}

resource "kubiya_agent" "product_agent" {
  name         = "product-knowledge-agent"
  runner       = "kubiya-hosted"
  description  = "Product team knowledge assistant"
  instructions = "You have access to product team Slack channels. Help with product decisions and feature discussions."
}
```

**Expected Outcome**: Creates separate knowledge integrations for different departments with specialized agents.

### 4. Incident Response Knowledge

Set up knowledge integration for incident response:

```hcl
resource "kubiya_external_knowledge" "incident_channels" {
  vendor = "slack"
  config = {
    channel_ids = [
      "C9999999999",  # #incidents
      "CAAAAAAAAAA",  # #incident-postmortems
      "CBBBBBBBBBB"   # #incident-logs
    ]
  }
}

resource "kubiya_agent" "incident_historian" {
  name         = "incident-historian"
  runner       = "kubiya-hosted"
  description  = "Incident response knowledge agent"
  instructions = <<-EOT
    You have access to incident-related Slack channels.
    Use this knowledge to:
    - Provide context on past incidents
    - Identify patterns in recurring issues
    - Reference previous solutions
    - Help with post-mortem analysis
    - Suggest preventive measures based on history
  EOT
  
  integrations = ["slack", "pagerduty"]
}
```

**Expected Outcome**: Creates an integration for incident response knowledge with an agent specialized in incident history.

### 5. Project-Based Knowledge

Configure knowledge for specific projects:

```hcl
locals {
  project_channels = {
    "project-alpha" = ["CCCCCCCCCCC", "CDDDDDDDDDD"]
    "project-beta"  = ["CEEEEEEEEE", "CFFFFFFFFFF"]
    "project-gamma" = ["CGGGGGGGGG", "CHHHHHHHHH"]
  }
}

resource "kubiya_external_knowledge" "project_knowledge" {
  for_each = local.project_channels
  
  vendor = "slack"
  config = {
    channel_ids = each.value
  }
}

resource "kubiya_agent" "project_agents" {
  for_each = local.project_channels
  
  name         = "${each.key}-agent"
  runner       = "kubiya-hosted"
  description  = "Knowledge agent for ${each.key}"
  instructions = "You have access to ${each.key} Slack channels. Help team members with project-specific questions and context."
}
```

**Expected Outcome**: Creates project-specific knowledge integrations with dedicated agents for each project.

### 6. Support Knowledge Base

Create a support knowledge integration:

```hcl
resource "kubiya_external_knowledge" "support_knowledge" {
  vendor = "slack"
  config = {
    channel_ids = [
      "CIIIIIIIIII",  # #customer-support
      "CJJJJJJJJJJ",  # #support-escalations
      "CKKKKKKKKKKK", # #support-internal
      "CLLLLLLLLLLL"  # #customer-feedback
    ]
  }
}

resource "kubiya_agent" "support_assistant" {
  name         = "support-knowledge-assistant"
  runner       = "kubiya-hosted"
  description  = "Customer support knowledge agent"
  instructions = <<-EOT
    You have access to support-related Slack channels.
    Use this knowledge to:
    - Find solutions to common customer issues
    - Reference past support interactions
    - Identify trending problems
    - Provide context on customer feedback
    - Help create support documentation
  EOT
  
  integrations = ["slack", "zendesk"]
}
```

**Expected Outcome**: Creates a comprehensive support knowledge base with an agent for customer support assistance.

### 7. Compliance and Security Knowledge

Set up knowledge for compliance tracking:

```hcl
resource "kubiya_external_knowledge" "compliance_knowledge" {
  vendor = "slack"
  config = {
    channel_ids = [
      "CMMMMMMMMMMM",  # #compliance
      "CNNNNNNNNNNN",  # #security-alerts
      "COOOOOOOOOO"    # #audit-logs
    ]
  }
}

resource "kubiya_agent" "compliance_agent" {
  name         = "compliance-knowledge-agent"
  runner       = "kubiya-hosted"
  description  = "Compliance and security knowledge agent"
  instructions = <<-EOT
    You have access to compliance and security Slack channels.
    IMPORTANT: This information is sensitive.
    
    Use this knowledge to:
    - Answer compliance-related questions
    - Provide context on security decisions
    - Reference audit discussions
    - Help with compliance reporting
    
    Always maintain confidentiality and follow security protocols.
  EOT
  
  groups = ["Compliance", "Security"]
}
```

**Expected Outcome**: Creates a secure knowledge integration for compliance and security teams.

## Argument Reference

### Required Arguments

* `vendor` - (Required, String) The vendor/provider for the knowledge integration. Currently supported values:
  - `slack` - Slack channel integrations

* `config` - (Required, Map of Dynamic) Dynamic configuration map with vendor-specific keys and values. The structure depends on the vendor being used.

### Vendor-Specific Configuration

#### Slack Configuration

For `vendor = "slack"`, the config map requires:

* `channel_ids` - (Required, List of Strings) List of Slack channel IDs to integrate. Each ID should be in the format "C1234567890".

Example:
```hcl
config = {
  channel_ids = ["C1234567890", "C0987654321"]
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the external knowledge integration.
* `org` - The organization associated with the integration.
* `start_date` - The start date from which knowledge is indexed.
* `integration_type` - The type of integration (matches the vendor).
* `created_at` - The timestamp when the integration was created.
* `updated_at` - The timestamp when the integration was last updated.

## Import

External knowledge integrations can be imported using their ID:

```shell
terraform import kubiya_external_knowledge.example <external-knowledge-id>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Slack integration requires OAuth setup in Kubiya dashboard
* Channel IDs must be valid and accessible
* Knowledge indexing may take time after initial creation
* Historical data retrieval limits may apply based on platform tier

## Best Practices

1. **Channel Selection**: Choose channels with valuable knowledge content
2. **Privacy Considerations**: Ensure compliance with data privacy policies
3. **Access Control**: Limit agent access to sensitive knowledge integrations
4. **Regular Updates**: Monitor and update channel lists as organizational structure changes
5. **Performance**: Avoid integrating extremely high-volume channels unnecessarily
6. **Testing**: Test knowledge retrieval in non-production environments first
7. **Documentation**: Document which channels are integrated and why
8. **Monitoring**: Monitor knowledge usage and relevance
9. **Retention Policies**: Align with your organization's data retention policies
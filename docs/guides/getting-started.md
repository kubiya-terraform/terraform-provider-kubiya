---
page_title: "Getting Started with Kubiya"
description: |-
  Getting started with the Kubiya provider
---

# Getting Started with Kubiya Provider

This guide will help you get started with the Kubiya Terraform provider. By the end of this guide, you'll have created your first Kubiya AI Teammate (agent) with a basic configuration and understand how to set up more advanced integrations.

## Prerequisites

- A Kubiya account with API access
- Terraform installed on your system
- Basic familiarity with Terraform

## Configuring the Provider

First, you'll need to configure the Kubiya provider in your Terraform configuration. Create a new directory for your Terraform files and create a file named `main.tf` with the following content:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key can be provided via environment variable KUBIYA_API_KEY
}
```

## Setting up your API Key

The Kubiya provider requires an API key for authentication. You can generate an API key from the Kubiya dashboard under Admin → Kubiya API Keys.

Set the API key as an environment variable:

```bash
export KUBIYA_API_KEY="your-api-key"
```

## Creating Your First Runner

Before you can create an agent, you need a runner. Runners are the compute environments where your agents execute code.

Add the following to your `main.tf` file:

```hcl
resource "kubiya_runner" "first_runner" {
  name        = "my-first-runner"
  description = "My first Kubiya runner"
  type        = "vcluster"
}
```

### Deploying the Runner

After creating the runner resource through Terraform, you will need to deploy the helm chart of the created runner. You can use the Kubiya dashboard to access the complete helm chart for your runner.

## Adding Knowledge Sources

Knowledge sources provide information to your agents. **Sources are required for agent resources** - every agent must have at least one source.

```hcl
# Create a source for your agent
resource "kubiya_source" "basic_tooling" {
  name        = "basic-tooling"
  description = "Basic tooling for the agent"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/basics"
  runner      = kubiya_runner.first_runner.name
}
```

## Creating Your First Agent (Sources Required)

Now you can create your first agent, which will use the runner and source you just defined. **Remember: Every agent must have at least one source** - this is a mandatory requirement:

```hcl
resource "kubiya_agent" "first_agent" {
  name        = "my-first-agent"
  runner      = kubiya_runner.first_runner.name
  description = "My first Kubiya agent created through Terraform"
  
  # Sources are REQUIRED for all agents - this is not optional
  sources = [kubiya_source.basic_tooling.name]
  
  # Optional: Configure some conversation starters
  starters = [
    {
      name    = "Hello"
      command = "Say hello and introduce yourself"
    },
    {
      name    = "Help"
      command = "List what you can help me with"
    }
  ]
}
```

## Adding Knowledge Resources

In addition to sources, you can create specific knowledge resources that contain information your agents can reference:

```hcl
# Add a direct knowledge entry
resource "kubiya_knowledge" "company_faq" {
  name        = "company-faq"
  description = "Company FAQ and procedures"
  content     = file("${path.module}/company-faq.md")
  format      = "markdown"
}
```

## Managing Sensitive Information with Secrets

Secrets allow you to securely store sensitive information that your agents might need:

```hcl
resource "kubiya_secret" "api_credentials" {
  name        = "api-credentials"
  description = "API credentials for external services"
  data = {
    api_key    = var.api_key
    api_secret = var.api_secret
  }
}

# Reference the secret in your agent
resource "kubiya_agent" "api_agent" {
  name        = "api-integration-agent"
  runner      = kubiya_runner.first_runner.name
  description = "Agent with API integration capabilities"
  
  # Sources are required for agents
  sources = [kubiya_source.basic_tooling.name]
  
  # Link the secret to the agent
  secrets = [kubiya_secret.api_credentials.name]
}
```

## Setting Up Integrations

Integrations connect your agents to external services:

```hcl
resource "kubiya_integration" "github_integration" {
  name        = "github-org"
  type        = "github"
  description = "GitHub organization integration"
  
  configuration = {
    org_name = "mycompany"
  }
}

# Link the integration to an agent
resource "kubiya_agent" "devops_agent" {
  name         = "devops-agent"
  runner       = kubiya_runner.first_runner.name
  description  = "DevOps automation agent"
  
  # Sources are required for agents
  sources      = [kubiya_source.basic_tooling.name]
  integrations = [kubiya_integration.github_integration.name]
}
```

## Creating Webhooks for Event-Driven Automation

Webhooks allow external systems to trigger your agents:

```hcl
resource "kubiya_webhook" "monitoring_webhook" {
  name        = "monitoring-alerts"
  agent       = kubiya_agent.devops_agent.name
  filter      = "alert.severity == 'critical'"
  source      = "prometheus"
  prompt      = "Investigate the critical alert from Prometheus"
  
  method      = "Slack"
  destination = "#alerts"
}
```

## Scheduling Recurring Tasks

Set up automated tasks that run on a schedule:

```hcl
resource "kubiya_scheduled_task" "daily_report" {
  name        = "daily-status-report"
  description = "Generate daily status report"
  agent       = kubiya_agent.devops_agent.name
  schedule    = "0 9 * * MON-FRI"  # 9 AM weekdays
  prompt      = "Generate a daily status report for our infrastructure"
  
  notification {
    method      = "Slack"
    destination = "#daily-reports"
  }
}
```

## Applying Your Configuration

Initialize your Terraform configuration:

```bash
terraform init
```

Plan your changes:

```bash
terraform plan
```

Apply your configuration:

```bash
terraform apply
```

When prompted, type `yes` to confirm the changes.

## Complete Basic Example

Here's a complete working example that brings together all the core resources we've discussed - runner, sources, knowledge, agent, and webhooks:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
  }
}

provider "kubiya" {
  # API key via KUBIYA_API_KEY environment variable
}

# 1. Create a runner
resource "kubiya_runner" "support_runner" {
  name        = "support-runner"
  description = "Runner for support automation"
  type        = "vcluster"
}

# 2. Create sources (REQUIRED for agents) - one for tooling and one for documentation
resource "kubiya_source" "support_tools" {
  name        = "support-tools"
  description = "Support automation tools"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/support"
  runner      = kubiya_runner.support_runner.name
}

resource "kubiya_source" "product_docs" {
  name        = "product-documentation"
  description = "Product documentation repository"
  url         = "https://github.com/mycompany/product-docs"
  branch      = "main"
  runner      = kubiya_runner.support_runner.name
}

# 3. Create knowledge resources for FAQs and troubleshooting guides
resource "kubiya_knowledge" "customer_faq" {
  name        = "customer-faq"
  description = "Frequently asked customer questions"
  content     = file("${path.module}/resources/faq.md")
  format      = "markdown"
}

resource "kubiya_knowledge" "troubleshooting" {
  name        = "troubleshooting-guide"
  description = "Troubleshooting procedures for common issues"
  content     = file("${path.module}/resources/troubleshooting.md")
  format      = "markdown"
  labels      = ["support", "troubleshooting"]
  
  # Specify which agent can use this knowledge
  supported_agents = ["support-assistant"]
}

# 4. Create a secret for API access
resource "kubiya_secret" "ticket_system" {
  name        = "ticket-system-api"
  description = "API credentials for the ticket system"
  data = {
    api_key    = var.ticket_system_api_key
    api_url    = "https://ticketing.example.com/api/v1"
    username   = "api-user"
  }
}

# 5. Create the support agent with sources, knowledge and secrets
resource "kubiya_agent" "support_assistant" {
  name        = "support-assistant"
  runner      = kubiya_runner.support_runner.name
  description = "AI-powered support assistant"
  
  # Sources are REQUIRED - agent must have at least one source
  sources = [
    kubiya_source.support_tools.name,
    kubiya_source.product_docs.name
  ]
  
  # Reference the secrets
  secrets = [kubiya_secret.ticket_system.name]
  
  # Add access control
  users  = ["support-team@example.com"]
  groups = ["Support"]
  
  # Conversation starters
  starters = [
    {
      name    = "Help Customer"
      command = "Help a customer with their issue"
    },
    {
      name    = "Troubleshoot"
      command = "Start troubleshooting a technical problem"
    }
  ]
  
  # Environment variables
  environment_variables = {
    DEFAULT_PRIORITY = "medium"
    LOG_LEVEL        = "info"
  }
}

# 6. Create a webhook for ticket system integration
resource "kubiya_webhook" "ticket_webhook" {
  name        = "new-ticket-alert"
  agent       = kubiya_agent.support_assistant.name
  source      = "TicketSystem"
  filter      = "ticket.priority == 'high'"
  prompt      = "A new high priority ticket has been created: {{.event.ticket.id}} - {{.event.ticket.subject}}"
  
  method      = "Slack"
  destination = "#support-alerts"
}

# 7. Create a scheduled task for daily report
resource "kubiya_scheduled_task" "daily_ticket_summary" {
  name        = "daily-ticket-summary"
  description = "Generate daily summary of support tickets"
  agent       = kubiya_agent.support_assistant.name
  schedule    = "0 17 * * MON-FRI"  # 5 PM weekdays
  prompt      = "Generate a summary of today's support tickets and response times"
  
  notification {
    method      = "Slack"
    destination = "#support-team"
  }
}

# Output important information
output "support_assistant" {
  value = {
    name            = kubiya_agent.support_assistant.name
    id              = kubiya_agent.support_assistant.id
    webhook_url     = kubiya_webhook.ticket_webhook.url
  }
}
```

This example demonstrates:
1. Creating a runner for compute
2. Creating multiple sources (required for agents)
3. Adding knowledge resources with content from files
4. Setting up a secret for secure API access
5. Configuring an agent with sources, secrets, and starters
6. Adding a webhook for event-driven automation
7. Setting up a scheduled task for recurring operations

It provides a complete foundation that you can customize for your specific use case.

## Advanced Example 1: CI/CD Maintainer Integration

The following is a more advanced example that sets up a CI/CD Maintainer agent integrated with GitHub for workflow monitoring and failure analysis:

```hcl
terraform {
  required_providers {
    kubiya = {
      source = "kubiya-terraform/kubiya"
    }
    github = {
      source  = "hashicorp/github"
      version = "6.4.0"
    }
  }
}

provider "kubiya" {
  # API key via KUBIYA_API_KEY environment variable
}

provider "github" {
  owner = local.github_organization
}

# Local variables for configuration
locals {
  repository_list = compact(split(",", var.repositories))
  github_events = ["check_run", "workflow_run"]
  webhook_filter = "workflow_run.conclusion != 'success' && workflow_run.event == 'pull_request'"
  github_organization = trim(split("/", local.repository_list[0])[0], " ")
}

# GitHub tools source
resource "kubiya_source" "github_tooling" {
  name = "github-tooling"
  url  = "https://github.com/kubiyabot/community-tools/tree/main/github"
}

# GitHub token secret
resource "kubiya_secret" "github_token" {
  name        = "GH_TOKEN"
  description = "GitHub token for the CI/CD Maintainer"
  data = {
    token = var.GITHUB_TOKEN
  }
}

# CI/CD Maintainer agent
resource "kubiya_agent" "cicd_maintainer" {
  name        = "cicd-maintainer"
  runner      = var.kubiya_runner
  description = "AI assistant that helps with GitHub Actions workflow failures"
  
  # Sources are required for agents
  sources      = [kubiya_source.github_tooling.name]
  secrets      = [kubiya_secret.github_token.name]
  integrations = ["slack"]
  
  environment_variables = {
    KUBIYA_TOOL_TIMEOUT = "500"
  }
}

# GitHub webhook
resource "kubiya_webhook" "github_webhook" {
  name        = "github-cicd-webhook"
  agent       = kubiya_agent.cicd_maintainer.name
  filter      = local.webhook_filter
  source      = "GitHub"
  method      = "Slack"
  destination = "#github-alerts"
  
  prompt      = <<-EOT
Your Goal: Analyze the failed GitHub Actions workflow.
Workflow ID: {{.event.workflow_run.id}}
PR Number: {{.event.workflow_run.pull_requests[0].number}}
Repository: {{.event.repository.full_name}}
  EOT
}

# Create GitHub repository webhooks
resource "github_repository_webhook" "repo_webhooks" {
  for_each = toset(local.repository_list)
  
  repository = try(trim(split("/", each.value)[1], " "), each.value)
  
  configuration {
    url          = kubiya_webhook.github_webhook.url
    content_type = "json"
    insecure_ssl = false
  }
  
  active = true
  events = local.github_events
}
```
## Next Steps

Now that you've created your first Kubiya resources with Terraform, you can:

1. Enhance your agents with more integrations and knowledge sources
2. Configure access control with users and groups
3. Build complex automation workflows with webhooks and scheduled tasks
4. Explore advanced use cases like the CI/CD Maintainer and JIT Permissions examples

Check out the resource documentation for detailed information on all the available resources and their configuration options. 
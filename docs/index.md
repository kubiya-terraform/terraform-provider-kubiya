---
page_title: "Provider: Kubiya"
description: |-
  Manage Kubiya AI resources, including agents, workflows, actions, and integrations, via the Kubiya API.
---

# Kubiya Provider

![Kubiya Logo](kubi.png)

The Kubiya provider is used to interact with the resources supported by the [Kubiya AI platform](https://kubiya.ai/). Manage Kubiya AI resources, including agents, workflows, actions, and integrations, via the Kubiya API. The provider needs to be configured with the proper credentials before it can be used.

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
  # Configuration options
}
```

## Authentication

The Kubiya provider uses environment variables for authentication. The API key is always read from the `KUBIYA_API_KEY` environment variable.

### Configuration

```hcl
provider "kubiya" {
  # API key is automatically read from KUBIYA_API_KEY environment variable
}
```

### Setting the Environment Variable

```sh
export KUBIYA_API_KEY="my-api-key"
```

To generate an API key, go to the Kubiya dashboard under Admin → Kubiya API Keys.

## Supported Resources

The following resources are supported by the Kubiya provider:

* [kubiya_agent](resources/agent.md) - Create and manage AI agents
* [kubiya_runner](resources/runner.md) - Configure agent execution environments
* [kubiya_integration](resources/integration.md) - Set up third-party integrations
* [kubiya_webhook](resources/webhook.md) - Configure webhooks for event-driven automation
* [kubiya_trigger](resources/trigger.md) - Create HTTP triggers for workflows
* [kubiya_secret](resources/secret.md) - Manage secure credentials
* [kubiya_source](resources/source.md) - Define tool and workflow sources
* [kubiya_knowledge](resources/knowledge.md) - Configure knowledge bases
* [kubiya_scheduled_task](resources/scheduled_task.md) - Set up scheduled automation tasks
* [kubiya_external_knowledge](resources/external_knowledge.md) - Connect external knowledge sources 
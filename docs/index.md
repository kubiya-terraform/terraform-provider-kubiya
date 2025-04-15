---
page_title: "Provider: Kubiya"
description: |-
  The Kubiya provider is used to interact with the Kubiya AI platform resources.
---

# Kubiya Provider

The Kubiya provider is used to interact with the resources supported by the [Kubiya AI platform](https://kubiya.ai/). The provider needs to be configured with the proper credentials before it can be used.

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

The Kubiya provider offers a flexible means of providing credentials for authentication. The following methods are supported, in this order:

- Static credentials
- Environment variables

### Static credentials

```hcl
provider "kubiya" {
  api_key = "my-api-key"
}
```

### Environment variables

```hcl
provider "kubiya" {}
```

```sh
export KUBIYA_API_KEY="my-api-key"
```

To generate an API key, go to the Kubiya dashboard under Admin → Kubiya API Keys.

## Supported Resources

The following resources are supported by the Kubiya provider:

* [kubiya_agent](resources/agent.md)
* [kubiya_runner](resources/runner.md)
* [kubiya_integration](resources/integration.md)
* [kubiya_webhook](resources/webhook.md)
* [kubiya_secret](resources/secret.md)
* [kubiya_source](resources/source.md)
* [kubiya_knowledge](resources/knowledge.md)
* [kubiya_scheduled_task](resources/scheduled_task.md) 
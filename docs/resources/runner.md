---
page_title: "kubiya_runner Resource - Kubiya"
subcategory: ""
description: |-
  The kubiya_runner resource manages compute environments for executing agents in the Kubiya platform.
---

# kubiya_runner (Resource)

The `kubiya_runner` resource allows you to create and manage runners in the Kubiya platform. Runners are compute environments where Kubiya agents execute code and perform operations. They provide isolated, secure execution contexts for your automation tasks.

## Prerequisites

Before using this resource, ensure you have:
1. A Kubiya account with API access
2. An API key (generated from Kubiya dashboard under Admin → Kubiya API Keys)
3. Appropriate permissions to create runners in your organization

## Example Usage

### 1. Basic Runner Configuration

Create a simple runner for general use:

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

resource "kubiya_runner" "basic" {
  name = "basic-runner"
}
```

**Expected Outcome**: Creates a basic runner that can be used by agents for task execution.

### 2. Production Runner

Configure a runner for production workloads:

```hcl
resource "kubiya_runner" "production" {
  name = "production-runner"
}

resource "kubiya_agent" "prod_agent" {
  name         = "production-agent"
  runner       = kubiya_runner.production.name
  description  = "Production automation agent"
  instructions = "You are a production agent running on dedicated infrastructure."
}
```

**Expected Outcome**: Creates a production-grade runner with an agent configured to use it.

### 3. Development Runner

Create a runner for development and testing:

```hcl
resource "kubiya_runner" "development" {
  name = "dev-runner"
}

resource "kubiya_agent" "dev_agent" {
  name         = "dev-assistant"
  runner       = kubiya_runner.development.name
  description  = "Development environment agent"
  instructions = "You are a development agent for testing and experimentation."
  
  environment_variables = {
    ENVIRONMENT = "development"
    DEBUG       = "true"
  }
}
```

**Expected Outcome**: Creates a development runner with an agent configured for testing.

### 4. Multi-Environment Setup

Create runners for different environments:

```hcl
# Define environments
locals {
  environments = ["dev", "staging", "prod"]
}

# Create runners for each environment
resource "kubiya_runner" "environment" {
  for_each = toset(local.environments)
  
  name = "${each.value}-runner"
}

# Create agents for each environment
resource "kubiya_agent" "environment_agent" {
  for_each = toset(local.environments)
  
  name         = "${each.key}-agent"
  runner       = kubiya_runner.environment[each.key].name
  description  = "Agent for ${each.key} environment"
  instructions = "You are an agent operating in the ${each.key} environment."
  
  environment_variables = {
    ENVIRONMENT = each.key
  }
}
```

**Expected Outcome**: Creates a complete multi-environment setup with dedicated runners and agents for each environment.

### 5. Specialized Task Runner

Create a runner for specific task types:

```hcl
resource "kubiya_runner" "data_processing" {
  name = "data-processing-runner"
}

resource "kubiya_source" "data_tools" {
  url    = "https://github.com/org/data-tools"
  runner = kubiya_runner.data_processing.name
}

resource "kubiya_agent" "data_agent" {
  name         = "data-processor"
  runner       = kubiya_runner.data_processing.name
  description  = "Agent specialized in data processing"
  instructions = "You are a data processing agent with access to pandas and data manipulation tools."
  
  sources = [kubiya_source.data_tools.id]
}
```

**Expected Outcome**: Creates a specialized runner with tools and an agent configured for data processing tasks.

## Argument Reference

### Required Arguments

* `name` - (Required, String) The name of the runner. Must be unique within your organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `runner_type` - The type of the runner (computed by the system).

## Import

Runners can be imported using their name:

```shell
terraform import kubiya_runner.example <runner-name>
```

## Compatibility Notes

* Requires Kubiya Terraform Provider version >= 1.0.0
* Compatible with Terraform >= 1.0
* Runners must be created before they can be referenced by agents or sources
* The special runner name "kubiya-hosted" refers to Kubiya's cloud-hosted runners

## Best Practices

1. **Environment Separation**: Create separate runners for different environments (dev, staging, prod)
2. **Naming Convention**: Use descriptive names that indicate the runner's purpose
3. **Resource Planning**: Consider workload requirements when planning runner deployment
4. **Agent Association**: Always specify appropriate runners for your agents based on their workload
5. **Monitoring**: Monitor runner status and performance for production workloads
6. **Documentation**: Document the purpose and configuration of each runner
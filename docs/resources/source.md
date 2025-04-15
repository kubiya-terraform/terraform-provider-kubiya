---
page_title: "kubiya_source Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Source resource to manage tool and workflow repositories.
---

# kubiya_source (Resource)

Provides a Kubiya Source resource. This allows sources to be created, updated, and deleted on the Kubiya platform. 

Sources in Kubiya are repositories containing tools and workflows that can be used by your AI Teammates (agents). **Sources are required for agents** - every agent must have at least one source connected to it.

When you add a source, Kubiya discovers the tools and workflows contained in the repository, making them available for use by your agents. Tools provide the functionality that agents can use to perform tasks, while workflows define sequences of operations.

## Container-First Architecture

Kubiya uses a container-first architecture where every tool is backed by a Docker image. This eliminates the need to write complex business logic from scratch and ensures secure, isolated execution environments for your tools.

Tools can be implemented in any language, with the Kubiya SDK providing first-class support for Python and JavaScript. When an agent invokes a tool, the execution happens in a containerized environment with predictable, secure results.

## Example Usage

```hcl
# GitHub repository source with tools
resource "kubiya_source" "kubernetes_tools" {
  name        = "kubernetes-tools"
  description = "Kubernetes management tools"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/kubernetes"
  runner      = kubiya_runner.main_runner.name
}

# Source with dynamic configuration for tool parameters
resource "kubiya_source" "aws_tools" {
  name        = "aws-tools"
  description = "AWS management tools"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/aws"
  runner      = kubiya_runner.main_runner.name
  
  # Configure default parameters for the tools
  dynamic_config = jsonencode({
    region          = "us-west-2"
    assume_role_arn = "arn:aws:iam::123456789012:role/KubiyaAssumeRole"
    environment     = "production"
  })
}

# Inline tools defined directly in Terraform
resource "kubiya_source" "inline_tools" {
  name        = "custom-tools"
  description = "Custom inline tools defined in Terraform"
  runner      = kubiya_runner.main_runner.name
  
  # Define custom tools directly in the source
  inline_tools = jsonencode([
    {
      name = "get_server_status"
      description = "Gets the status of a server"
      parameters = {
        server_name = {
          type = "string"
          description = "The name of the server to check"
        }
      }
      code = <<-EOT
        async function run(params) {
          const { server_name } = params;
          // Custom logic to check server status
          return { status: "online", server: server_name };
        }
      EOT
    },
    {
      name = "restart_service"
      description = "Restarts a service on the specified server"
      parameters = {
        server_name = {
          type = "string"
          description = "The name of the server"
        },
        service_name = {
          type = "string"
          description = "The name of the service to restart"
        }
      }
      code = <<-EOT
        async function run(params) {
          const { server_name, service_name } = params;
          // Custom logic to restart service
          return { result: "service restarted", server: server_name, service: service_name };
        }
      EOT
    }
  ])
  
  # Define custom workflows that use the tools
  inline_workflows = jsonencode([
    {
      name = "check_and_restart"
      description = "Checks server status and restarts service if needed"
      steps = [
        {
          tool = "get_server_status"
          params = {
            server_name = "{{.server_name}}"
          }
          save_result_as = "status_result"
        },
        {
          condition = "{{.status_result.status}} == 'online'"
          tool = "restart_service"
          params = {
            server_name = "{{.server_name}}",
            service_name = "{{.service_name}}"
          }
        }
      ]
    }
  ])
}

# Private repository with authentication
resource "kubiya_source" "private_tools" {
  name        = "private-tools"
  description = "Private internal tools"
  url         = "https://github.com/mycompany/private-tools"
  branch      = "main"
  runner      = kubiya_runner.main_runner.name
  auth_method = "ssh_key"
}

# Python SDK tool example
resource "kubiya_source" "python_tools" {
  name        = "python-tools"
  description = "Tools created with the Kubiya Python SDK"
  url         = "https://github.com/mycompany/kubiya-python-tools"
  runner      = kubiya_runner.main_runner.name
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the source.
* `description` - (Required) A description of the source.
* `runner` - (Required) The runner that will be used for tool execution.

### Optional Arguments

* `url` - (Optional) The URL of the repository containing tools and workflows. Required unless using inline tools.
* `branch` - (Optional) For Git sources, the branch to use. Defaults to "main" or "master" depending on the repository.
* `auth_method` - (Optional) Authentication method for private sources. Values: "ssh_key", "token", etc.
* `dynamic_config` - (Optional) JSON-encoded configuration data to pass to the tools in the source. This is often used to configure tools with specific parameters or credentials.
* `inline_tools` - (Optional) JSON-encoded array of tool definitions to be created directly in the source instead of from a repository. Each tool definition includes name, description, parameters, and code.
* `inline_workflows` - (Optional) JSON-encoded array of workflow definitions to be created directly in the source. Each workflow definition includes name, description, and a series of steps that use tools.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the source.
* `status` - The current status of the source.
* `last_updated` - The timestamp when the source was last updated or synchronized.

## Tools and Workflows

Sources in Kubiya provide two main types of components:

1. **Tools**: Functions that agents can use to perform specific tasks, such as interacting with AWS, managing Kubernetes clusters, or querying databases. Tools can be defined in a repository or inline.

2. **Workflows**: Sequences of tool calls that define a process. Workflows allow chaining tools together with conditional logic and variable passing between steps.

When Kubiya processes a source, it discovers all tools and workflows available in the repository or defined inline, making them available to any agent that has that source attached.

## Creating Tools with the Kubiya SDK

While you can define tools directly in Terraform using the `inline_tools` parameter, Kubiya also offers a Python SDK for creating more complex tools. The SDK provides a way to define tools that will run in a containerized environment:

```python
from kubiya_sdk import tool

@tool(image="python:3.12-slim")
def hello_world(name: str) -> str:
    return f"Hello, {name}!"

# More complex example
@tool(name="database_migration")
def migrate_database(source_db: str, target_db: str, tables: list) -> dict:
    """Execute a secure database migration with validation"""
    # Implementation runs in containerized environment
    return {
        "status": "success", 
        "tables_migrated": len(tables),
        "validation_passed": True
    }
```

These tools can be packaged into a repository and then referenced by the `url` parameter in your `kubiya_source` resource. The Kubiya platform provides a growing repository of pre-built tools for common enterprise tasks at [github.com/kubiyabot/community-tools](https://github.com/kubiyabot/community-tools).

## Dynamic Configuration

The `dynamic_config` parameter allows you to pass configuration data to the tools in a source. This is useful for setting default parameters, credentials, or environment-specific values that the tools will use during execution.

For example, if your source contains AWS tools, you might want to set the default region and role to use:

```hcl
dynamic_config = jsonencode({
  region = "us-west-2"
  role_arn = "arn:aws:iam::123456789012:role/KubiyaRole"
})
```

## Import

Sources can be imported using the `id`:

```
$ terraform import kubiya_source.example SOURCE_ID
```

## Connecting Sources to Agents

After creating a source, you need to connect it to an agent to make its tools and workflows available. This is done using the `sources` attribute in the agent resource:

```hcl
resource "kubiya_agent" "example" {
  name        = "example-agent"
  runner      = kubiya_runner.example.name
  description = "Example agent with tools"
  
  # Connect sources to the agent
  sources = [
    kubiya_source.kubernetes_tools.name,
    kubiya_source.aws_tools.name,
    kubiya_source.inline_tools.name,
    kubiya_source.python_tools.name
  ]
}
```

Remember that **every agent must have at least one source** - this is a mandatory requirement. 
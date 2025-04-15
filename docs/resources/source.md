---
page_title: "kubiya_source Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Source resource to manage tool and workflow repositories.
---

# kubiya_source (Resource)

Provides a Kubiya Source resource. This allows sources to be created, updated, and deleted on the Kubiya platform. 

Sources in Kubiya are repositories containing tools and workflows that can be used by your AI Teammates (agents). **Sources are required for agents** - every agent must have at least one source connected to it.

When you add a source, Kubiya discovers the tools and workflows contained in the repository, making them available for use by your agents.

## Example Usage

```hcl
# GitHub repository source
resource "kubiya_source" "kubernetes_tools" {
  name        = "kubernetes-tools"
  description = "Kubernetes management tools"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/kubernetes"
  runner      = kubiya_runner.main_runner.name
}

# GitLab repository source
resource "kubiya_source" "internal_tools" {
  name        = "internal-tools"
  description = "Internal company tools"
  url         = "https://gitlab.com/mycompany/tools"
  branch      = "master"
  runner      = kubiya_runner.main_runner.name
}

# Source with dynamic configuration
resource "kubiya_source" "aws_tools" {
  name        = "aws-tools"
  description = "AWS management tools"
  url         = "https://github.com/kubiyabot/community-tools/tree/main/aws"
  runner      = kubiya_runner.main_runner.name
  dynamic_config = jsonencode({
    region    = "us-west-2"
    role_arn  = "arn:aws:iam::123456789012:role/KubiyaAssumeRole"
  })
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
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the source.
* `description` - (Required) A description of the source.
* `url` - (Required) The URL of the repository containing tools and workflows.
* `runner` - (Required) The runner that will be used for tool execution.

### Optional Arguments

* `branch` - (Optional) For Git sources, the branch to use. Defaults to "main" or "master" depending on the repository.
* `auth_method` - (Optional) Authentication method for private sources. Values: "ssh_key", "token", etc.
* `dynamic_config` - (Optional) JSON-encoded configuration data to pass to the tools in the source. This is often used to configure tools with specific parameters or credentials.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the source.
* `status` - The current status of the source.
* `last_updated` - The timestamp when the source was last updated or synchronized.

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
    kubiya_source.aws_tools.name
  ]
}
```

Remember that **every agent must have at least one source** - this is a mandatory requirement. 
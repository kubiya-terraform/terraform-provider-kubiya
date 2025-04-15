---
page_title: "kubiya_agent Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Agent resource to manage AI Teammates on the Kubiya platform.
---

# kubiya_agent (Resource)

Provides a Kubiya Agent resource. This allows AI Teammates to be created, updated, and deleted on the Kubiya platform.

## Example Usage

```hcl
resource "kubiya_agent" "example" {
  name        = "kubernetes-assistant"
  runner      = "runnerv2-5-vcluster"
  description = "This agent can perform various Kubernetes tasks using kubectl"
  
  # Optional fields with defaults
  model        = "gpt-4o"  # Defaults to "gpt-4o"
  is_debug_mode = false
  
  # Access control
  users  = ["user@example.com"]
  groups = ["Administrators"]
  
  # Integrations and capabilities
  integrations = ["Github"]
  sources      = ["kubernetes-docs"]
  secrets      = ["kube-config"]
  
  # Environment variables
  environment_variables = {
    KUBE_NAMESPACE = "default"
  }
  
  # Conversation starters
  starters = [
    {
      name    = "List Pods"
      command = "kubectl get pods"
    },
    {
      name    = "Check Node Status"
      command = "kubectl get nodes"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the agent.
* `runner` - (Required) The runner that will execute the agent. Must reference an existing runner.
* `description` - (Required) A detailed description of the agent's capabilities and purpose.

### Optional Arguments with Defaults

* `model` - (Optional) The LLM model to use. Defaults to "gpt-4o".
* `is_debug_mode` - (Optional) Whether to enable debug mode. Defaults to false.

### Optional Arguments

* `users` - (Optional) List of users that have access to this agent. Each user must exist in the Kubiya platform.
* `groups` - (Optional) List of groups that have access to this agent. Each group must exist in the Kubiya platform.
* `sources` - (Optional) List of knowledge sources to be used by the agent. Each source must exist in the Kubiya platform.
* `secrets` - (Optional) List of secrets accessible to the agent. Each secret must exist in the Kubiya platform.
* `integrations` - (Optional) List of integrations to enable for the agent. Each integration must exist in the Kubiya platform.
* `environment_variables` - (Optional) Map of environment variables for the agent.

* `starters` - (Optional) List of conversation starters. Each starter has:
  * `name` - The display name of the starter.
  * `command` - The command or prompt for the starter.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the agent.
* `owner` - The owner of the agent.
* `created_at` - The timestamp when the agent was created.

## Import

Agents can be imported using the `id`:

```
$ terraform import kubiya_agent.example AGENT_ID
``` 
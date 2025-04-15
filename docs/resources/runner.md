---
page_title: "kubiya_runner Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Runner resource to manage compute environments for agents.
---

# kubiya_runner (Resource)

Provides a Kubiya Runner resource. This allows runners to be created, updated, and deleted. Runners are compute environments where Kubiya agents execute code and perform operations.

## Example Usage

```hcl
# Example of a standard runner configuration
resource "kubiya_runner" "standard" {
  name        = "production-runner"
  description = "Production runner for critical operations"
}

```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the runner.
* `description` - (Required) A description of the runner's purpose.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the runner.
* `status` - The current status of the runner.
* `created_at` - The timestamp when the runner was created.

## Import

Runners can be imported using the `id`:

```
$ terraform import kubiya_runner.example RUNNER_ID
``` 
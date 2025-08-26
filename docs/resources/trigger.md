---
page_title: "kubiya_trigger Resource - terraform-provider-kubiya"
subcategory: ""
description: |-
  Manages a Kubiya workflow trigger with webhook capabilities
---

# kubiya_trigger (Resource)

Manages a Kubiya workflow trigger with webhook capabilities. This resource creates a workflow and automatically
publishes it with a webhook trigger, providing a URL that can be used to execute the workflow.

## Example Usage

### Simple Echo Workflow

```terraform
resource "kubiya_trigger" "workflow_trigger" {
  name   = "workflow-trigger"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Echo via Webhook"
    version = 1
    steps = [
      {
        name = "echo"
        executor = {
          type = "command"
          config = {
            command = "echo \"Hello from webhook\""
          }
        }
      }
    ]
  })
}

output "trigger_url" {
  value = kubiya_trigger.workflow_trigger.url
}
```

### Multi-Step Workflow

```terraform
resource "kubiya_trigger" "multi_step" {
  name   = "data-processing"
  runner = "kubiya-hosted"
  
  workflow = jsonencode({
    name    = "Data Processing Pipeline"
    version = 1
    steps = [
      {
        name = "fetch_data"
        executor = {
          type = "command"
          config = {
            command = "curl -s https://api.example.com/data"
          }
        }
      },
      {
        name = "process_data"
        executor = {
          type = "command"
          config = {
            command = "jq '.items | length'"
          }
        }
      },
      {
        name = "notify"
        executor = {
          type = "command"
          config = {
            command = "echo \"Processing complete\""
          }
        }
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the trigger. This will be used as the workflow name.
* `runner` - (Required) Runner to use for executing the workflow. Common values include:
    - `kubiya-hosted` - Use Kubiya's hosted runners
    - Custom runner names from your organization
* `workflow` - (Required) JSON-encoded workflow definition (use `jsonencode()`) containing:
    - `name` - (Required) Name of the workflow
    - `version` - (Required) Version number of the workflow (integer)
    - `steps` - (Required) List of workflow steps, each containing:
        - `name` - (Required) Name of the step
        - `executor` - (Required) Executor configuration:
            - `type` - (Required) Type of executor (e.g., "command")
            - `config` - (Required) Configuration for the executor:
                - `command` - (Required) Command to execute

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the trigger
* `url` - The webhook URL for triggering the workflow. POST requests to this URL will execute the workflow.
* `status` - Current status of the workflow (e.g., "draft", "published")
* `workflow_id` - The ID of the created workflow in Kubiya

## Triggering the Workflow

Once the trigger resource is created, you can execute the workflow by making a POST request to the webhook URL:

```bash
# Get the webhook URL from Terraform output
WEBHOOK_URL=$(terraform output -raw trigger_url)

# Trigger the workflow
curl -X POST "$WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: UserKey YOUR_API_KEY" \
  -d '{"key": "value"}'
```

For streaming execution output, append `?stream=true` to the webhook URL:

```bash
curl -X POST "$WEBHOOK_URL?stream=true" \
  -H "Content-Type: application/json" \
  -H "Authorization: UserKey YOUR_API_KEY" \
  -d '{}'
```

## Import

Trigger resources can be imported using their ID:

```shell
terraform import kubiya_trigger.example <trigger-id>
```

## Notes

- The workflow is automatically published when the trigger resource is created
- The webhook URL remains stable across updates unless the resource is recreated
- Updating the workflow definition will update the published workflow
- Deleting the trigger resource will delete both the workflow and webhook
---
page_title: "kubiya_webhook Resource - terraform-provider-kubiya"
description: |-
  Provides a Kubiya Webhook resource to manage webhook triggers for agents.
---

# kubiya_webhook (Resource)

Provides a Kubiya Webhook resource. This allows webhooks to be created, updated, and deleted on the Kubiya platform, enabling external systems to trigger Kubiya agents via HTTP requests.

## Example Usage

```hcl
resource "kubiya_webhook" "example" {
  name        = "github-issues"
  agent       = kubiya_agent.example.name
  filter      = "issues"
  source      = "github"
  prompt      = "Process this new GitHub issue"
  
  # Notification settings
  method      = "Slack"
  destination = "#notifications"
}

# Example with Teams notification
resource "kubiya_webhook" "teams_example" {
  name        = "deployment-notification"
  agent       = kubiya_agent.example.name
  source      = "github-actions"
  
  method      = "Teams"
  team_name   = "DevOps Team"
  destination = "Deployments"  # Channel name in Teams
}

# Example with HTTP notification (no destination required)
resource "kubiya_webhook" "http_example" {
  name        = "monitoring-alert"
  agent       = kubiya_agent.example.name
  source      = "prometheus"
  
  method      = "Http"
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the webhook.
* `agent` - (Required) The name of the agent to trigger with this webhook.

### Optional Arguments

* `filter` - (Optional) Filter for events that will trigger the webhook.
* `source` - (Optional) Source identification for the webhook.
* `prompt` - (Optional) Prompt to send to the agent when the webhook is triggered.

* `method` - (Optional) Notification method. Values: "Slack", "Teams", "Http". Defaults to "Slack".
* `destination` - (Optional) Destination for notifications:
  * For Slack: Channel name with "#" prefix (e.g., "#alerts", not needed for MS Teams) or username with "@" prefix
  * For MS Teams: Channel name within the team specified by `team_name`
  * For Http: Not required
  
* `team_name` - (Optional) Team name for Microsoft Teams notifications. Required when `method` is "Teams".

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the webhook.
* `url` - The URL to trigger the webhook.
* `created_at` - The timestamp when the webhook was created.
* `created_by` - The user who created the webhook.

## Import

Webhooks can be imported using the `id`:

```
$ terraform import kubiya_webhook.example WEBHOOK_ID
``` 
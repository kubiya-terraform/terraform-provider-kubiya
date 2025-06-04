---
page_title: "kubiya_external_knowledge Resource - terraform-provider-kubiya"
subcategory: ""
description: |-
  Manages external knowledge integrations in Kubiya.
---

# kubiya_external_knowledge (Resource)

The `kubiya_external_knowledge` resource manages external knowledge integrations in Kubiya. This resource allows you to connect external systems to the Kubiya RAG (Retrieval-Augmented Generation) system through a vendor-agnostic interface.

## Example Usage

### Basic Slack Integration

```hcl
# Slack knowledge integration with single channel
resource "kubiya_external_knowledge" "slack_single" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890"]  # List with single channel ID
  }
}

# Slack with multiple channels
resource "kubiya_external_knowledge" "slack_multiple" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890", "C0987654321", "C1111111111"]
  }
}
```

## Schema

### Required

- `vendor` (String) - The vendor/provider for the knowledge integration. Currently supported:
  * `"slack"` - Slack channel integrations

- `config` (Map of Dynamic) - Dynamic configuration map with vendor-specific keys and values. The structure depends on the vendor being used.

### Read-Only

- `id` (String) - The unique identifier of the external knowledge integration
- `org` (String) - The organization associated with the integration
- `start_date` (String) - The start date of the integration
- `integration_type` (String) - The type of integration (matches the vendor)
- `created_at` (String) - The timestamp when the integration was created
- `updated_at` (String) - The timestamp when the integration was last updated

## Vendor-Specific Configuration

### Slack (`vendor = "slack"`)

For Slack integrations, the `config` map supports:

* `channel_ids` (List of Strings, Required) - List of Slack channel IDs to integrate
  * Always use a list, even for single channels
  * Example: `["C1234567890", "C0987654321"]`

Example:
```hcl
resource "kubiya_external_knowledge" "slack_channels" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890", "C0987654321"]
  }
}
```

## Import

External knowledge integrations can be imported using the vendor and ID separated by a comma:

```shell
terraform import kubiya_external_knowledge.example slack,550e8400-e29b-41d4-a716-446655440000
```

## Migration from Legacy Resources

If you're migrating from the old `kubiya_slack_knowledge` resource:

### Before (Old)
```hcl
resource "kubiya_slack_knowledge" "example" {
  channel_id = "C1234567890"
}
```

### After (New)
```hcl
resource "kubiya_external_knowledge" "example" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890"]  # Note: Always use a list for Slack
  }
}
```

## Finding Vendor-Specific IDs

### Slack Channel IDs
1. Open Slack in your web browser
2. Navigate to the channel you want to integrate
3. Click on the channel name at the top of the screen
4. At the bottom of the popup, you'll see the Channel ID (starts with 'C')
5. Alternatively, right-click on the channel name in the sidebar and select "Copy link" - the ID is in the URL

## Notes

- The `config` field is dynamic and can accept different data types (strings, lists, booleans, numbers) depending on the vendor's requirements
- Each vendor may have different required and optional fields in the config
- Always refer to the vendor-specific documentation section for the correct configuration structure
- When using lists in the config, ensure they are properly formatted as Terraform lists 
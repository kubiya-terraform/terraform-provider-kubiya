# External Knowledge Resource Examples

This directory contains examples of using the `kubiya_external_knowledge` resource to integrate external knowledge sources with Kubiya.

## Overview

The `kubiya_external_knowledge` resource allows you to connect various external knowledge sources to Kubiya. Currently, the following vendor is supported:

- **Slack** - Index and search Slack channels

## Prerequisites

1. Set your Kubiya API key:
   ```bash
   export KUBIYA_API_KEY="your-api-key-here"
   ```

2. Ensure you have the necessary permissions and credentials for the external service you want to integrate.

## Slack Integration

### Basic Example - Single Channel

```hcl
resource "kubiya_external_knowledge" "slack_single_channel" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890"]  # Replace with your Slack channel ID
  }
}
```

### Multiple Channels

```hcl
resource "kubiya_external_knowledge" "slack_multiple_channels" {
  vendor = "slack"
  config = {
    channel_ids = ["C1234567890", "C0987654321", "C1111111111"]
  }
}
```

### Finding Slack Channel IDs

To find a Slack channel ID:
1. Open Slack in your web browser
2. Navigate to the channel
3. Click on the channel name at the top
4. At the bottom of the popup, you'll see the Channel ID (starts with 'C')

## Resource Arguments

### Common Arguments

- `vendor` (Required) - The vendor/integration type. Currently supports: `"slack"`
- `config` (Required) - A map of configuration values specific to the vendor

### Slack Configuration

For Slack integrations, the `config` map must include:
- `channel_ids` (Required) - A list of Slack channel IDs to index

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The unique identifier of the external knowledge integration
- `org` - The organization ID
- `start_date` - When the integration started indexing
- `integration_type` - The type of integration (matches vendor)
- `created_at` - When the integration was created
- `updated_at` - When the integration was last updated

## Using Data Sources

You can reference existing external knowledge integrations:

```hcl
data "kubiya_external_knowledge" "example" {
  id     = kubiya_external_knowledge.slack_single_channel.id
  vendor = "slack"
}
```

## Complete Example

See `main.tf` for a complete working example that includes:
- Creating Slack integrations
- Using data sources
- Outputting integration details

## Running the Example

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Review the plan:
   ```bash
   terraform plan
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

## Adding New Vendors

To add support for new vendors, you need to:
1. Implement the vendor-specific logic in the provider code
2. Register the vendor in the vendor registry
3. Update the documentation

See the provider's `internal/clients/vendors/README.md` for implementation details 
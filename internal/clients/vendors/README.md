# External Knowledge Vendor Architecture

This directory contains the vendor-specific implementations for the `external_knowledge` resource in the Kubiya Terraform provider.

## Architecture Overview

The vendor architecture uses an interface-based design pattern that allows easy addition of new vendors without modifying the core client code.

### Key Components

1. **VendorClient Interface** (`interface.go`)
   - Defines the contract that all vendor implementations must satisfy
   - Provides a registry for managing vendor implementations

2. **Base Types** (`base.go`)
   - Common structures and helper functions used by all vendors
   - `BaseExternalKnowledge`: Common fields for all vendor responses
   - Helper functions for converting between Terraform and Go types

3. **Vendor Implementations**
   - `slack.go`: Slack-specific implementation

## Adding a New Vendor

To add support for a new vendor:

1. **Create a new vendor file** (e.g., `notion.go`)
   ```go
   package vendors
   
   type NotionVendor struct {
       name string
   }
   
   func NewNotionVendor() VendorClient {
       return &NotionVendor{name: "notion"}
   }
   ```

2. **Implement the VendorClient interface**
   - `GetVendorName()`: Return the vendor identifier
   - `ValidateConfig()`: Validate vendor-specific configuration
   - `PrepareCreateRequest()`: Convert Terraform config to API request
   - `PrepareUpdateRequest()`: Convert Terraform config to update request
   - `ParseCreateResponse()`: Parse API create response
   - `ParseReadResponse()`: Parse API read response
   - `ParseUpdateResponse()`: Parse API update response
   - `ParseListResponse()`: Parse API list response

3. **Register the vendor** in `interface.go`
   ```go
   func InitializeRegistry() *Registry {
       registry := NewRegistry()
       registry.Register(NewSlackVendor())
       registry.Register(NewNotionVendor()) // Add this line
       return registry
   }
   ```

## Vendor-Specific Request/Response Handling

Each vendor can define its own:
- Request structures (what gets sent to the API)
- Response structures (what comes back from the API)
- Configuration validation rules
- Field mappings between Terraform and the API

### Example: Slack Implementation

Slack requires `channel_ids` as a list of strings:

```go
type slackIntegrationRequest struct {
    ChannelIDs []string `json:"channel_ids"`
}
```

The Slack vendor validates that `channel_ids` is present and non-empty:

```go
func (s *SlackVendor) ValidateConfig(config types.Map) error {
    // Check if channel_ids exists and is not empty
    if v, ok := config.Elements()["channel_ids"]; ok {
        // Validation logic
    }
    return fmt.Errorf("channel_ids is required for Slack integration")
}
```

## Important Notes

- **Only explicitly registered vendors are supported** - If a vendor is not registered in `InitializeRegistry()`, it will return an error
- **Each vendor must have its own implementation** - There is no generic fallback

## Testing

When adding a new vendor:
1. Test the configuration validation
2. Test request preparation with various input types
3. Test response parsing for all CRUD operations
4. Add example Terraform configurations in `examples/external_knowledge/` 
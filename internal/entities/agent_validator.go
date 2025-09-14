package entities

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ValidateMCPServer validates the MCP server configuration
// This function is exported and used by the provider package
func ValidateMCPServer(mcp *MCPServerModel) error {
	if mcp == nil {
		return nil
	}

	serverType := mcp.Type.ValueString()

	switch serverType {
	case "stdio":
		// For STDIO type, command is required
		if mcp.Command.IsNull() || mcp.Command.IsUnknown() || mcp.Command.ValueString() == "" {
			return fmt.Errorf("command is required for stdio type MCP server")
		}

		// SSE fields should not be set for STDIO type
		if !mcp.URL.IsNull() && !mcp.URL.IsUnknown() && mcp.URL.ValueString() != "" {
			return fmt.Errorf("url should not be set for stdio type MCP server")
		}
		if !mcp.Headers.IsNull() && !mcp.Headers.IsUnknown() && len(mcp.Headers.Elements()) > 0 {
			return fmt.Errorf("headers should not be set for stdio type MCP server")
		}

	case "sse":
		// For SSE type, URL is required
		if mcp.URL.IsNull() || mcp.URL.IsUnknown() || mcp.URL.ValueString() == "" {
			return fmt.Errorf("url is required for sse type MCP server")
		}

		// STDIO fields should not be set for SSE type
		if !mcp.Command.IsNull() && !mcp.Command.IsUnknown() && mcp.Command.ValueString() != "" {
			return fmt.Errorf("command should not be set for sse type MCP server")
		}
		if !mcp.Args.IsNull() && !mcp.Args.IsUnknown() && len(mcp.Args.Elements()) > 0 {
			return fmt.Errorf("args should not be set for sse type MCP server")
		}
		if !mcp.Env.IsNull() && !mcp.Env.IsUnknown() && len(mcp.Env.Elements()) > 0 {
			return fmt.Errorf("env should not be set for sse type MCP server")
		}

	case "":
		return fmt.Errorf("type is required for MCP server configuration")

	default:
		return fmt.Errorf("invalid MCP server type '%s', must be 'stdio' or 'sse'", serverType)
	}

	return nil
}

// ValidateAgent validates the entire agent configuration
// This can be extended to validate other agent fields in the future
func ValidateAgent(agent *AgentModel) error {
	// Validate MCP server if present
	if agent.MCPServer != nil {
		if err := ValidateMCPServer(agent.MCPServer); err != nil {
			return fmt.Errorf("MCP server validation failed: %w", err)
		}
	}

	// Add other agent-level validations here in the future
	// For example:
	// - Validate runner exists
	// - Validate model is supported
	// - Validate environment variables format
	// - etc.

	return nil
}

// MCPServerTypeValidator validates the MCP server type field
// This validator is used in the schema definition
type MCPServerTypeValidator struct{}

func (v MCPServerTypeValidator) Description(ctx context.Context) string {
	return "MCP server type must be 'stdio' or 'sse'"
}

func (v MCPServerTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "MCP server type must be either `stdio` or `sse`"
}

func (v MCPServerTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value != "stdio" && value != "sse" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MCP Server Type",
			fmt.Sprintf("MCP server type must be 'stdio' or 'sse', got '%s'", value),
		)
	}
}

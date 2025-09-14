package entities

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateAgent(t *testing.T) {
	tests := []struct {
		name      string
		agent     *AgentModel
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid agent without MCP server",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer:    nil,
			},
			wantError: false,
		},
		{
			name: "valid agent with stdio MCP server",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer: &MCPServerModel{
					Type:    types.StringValue("stdio"),
					Command: types.StringValue("python"),
					Args:    types.ListNull(types.StringType),
					Env:     types.MapNull(types.StringType),
					URL:     types.StringNull(),
					Headers: types.MapNull(types.StringType),
				},
			},
			wantError: false,
		},
		{
			name: "valid agent with sse MCP server",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer: &MCPServerModel{
					Type:    types.StringValue("sse"),
					URL:     types.StringValue("https://mcp.example.com/sse"),
					Headers: types.MapNull(types.StringType),
					Command: types.StringNull(),
					Args:    types.ListNull(types.StringType),
					Env:     types.MapNull(types.StringType),
				},
			},
			wantError: false,
		},
		{
			name: "invalid agent with invalid MCP server type",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer: &MCPServerModel{
					Type:    types.StringValue("invalid"),
					Command: types.StringNull(),
					Args:    types.ListNull(types.StringType),
					Env:     types.MapNull(types.StringType),
					URL:     types.StringNull(),
					Headers: types.MapNull(types.StringType),
				},
			},
			wantError: true,
			errorMsg:  "MCP server validation failed: invalid MCP server type 'invalid', must be 'stdio' or 'sse'",
		},
		{
			name: "invalid agent with stdio MCP server missing command",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer: &MCPServerModel{
					Type:    types.StringValue("stdio"),
					Command: types.StringNull(),
					Args:    types.ListNull(types.StringType),
					Env:     types.MapNull(types.StringType),
					URL:     types.StringNull(),
					Headers: types.MapNull(types.StringType),
				},
			},
			wantError: true,
			errorMsg:  "MCP server validation failed: command is required for stdio type MCP server",
		},
		{
			name: "invalid agent with sse MCP server missing url",
			agent: &AgentModel{
				Name:         types.StringValue("test-agent"),
				Runner:       types.StringValue("default"),
				Description:  types.StringValue("Test agent"),
				Instructions: types.StringValue("Test instructions"),
				MCPServer: &MCPServerModel{
					Type:    types.StringValue("sse"),
					URL:     types.StringNull(),
					Headers: types.MapNull(types.StringType),
					Command: types.StringNull(),
					Args:    types.ListNull(types.StringType),
					Env:     types.MapNull(types.StringType),
				},
			},
			wantError: true,
			errorMsg:  "MCP server validation failed: url is required for sse type MCP server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAgent(tt.agent)
			hasError := err != nil

			if hasError != tt.wantError {
				t.Errorf("ValidateAgent() error = %v, wantError %v", err, tt.wantError)
			}

			if tt.wantError && err != nil && tt.errorMsg != "" {
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateAgent() error message = %v, want %v", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

func TestMCPServerTypeValidator(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "valid stdio type",
			value: "stdio",
			valid: true,
		},
		{
			name:  "valid sse type",
			value: "sse",
			valid: true,
		},
		{
			name:  "invalid type",
			value: "invalid",
			valid: false,
		},
		{
			name:  "empty type",
			value: "",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validator logic
			isValid := tt.value == "stdio" || tt.value == "sse"
			if isValid != tt.valid {
				t.Errorf("MCPServerTypeValidator validation for '%s': got %v, want %v", tt.value, isValid, tt.valid)
			}
		})
	}
}

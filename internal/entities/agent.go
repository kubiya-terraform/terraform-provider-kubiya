package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TaskModel struct {
	Name        string `tfsdk:"name"`
	Prompt      string `tfsdk:"prompt"`
	Description string `tfsdk:"description"`
}

type AgentModel struct {
	// Computed
	Id        types.String `tfsdk:"id"`
	Owner     types.String `tfsdk:"owner"`
	CreatedAt types.String `tfsdk:"created_at"`

	// Required with default value
	Image       types.String `tfsdk:"image"`
	Model       types.String `tfsdk:"model"`
	IsDebugMode types.Bool   `tfsdk:"is_debug_mode"`

	// Required
	Name         types.String `tfsdk:"name"`
	Runner       types.String `tfsdk:"runner"`
	Description  types.String `tfsdk:"description"`
	Instructions types.String `tfsdk:"instructions"`

	// Optional
	Links        types.List     `tfsdk:"links"`
	Tasks        []TaskModel    `tfsdk:"tasks"`
	Users        types.List     `tfsdk:"users"`
	Groups       types.List     `tfsdk:"groups"`
	Sources      types.List     `tfsdk:"sources"`
	Secrets      types.List     `tfsdk:"secrets"`
	Starters     []StarterModel `tfsdk:"starters"`
	Workflows    types.String   `tfsdk:"workflows"`
	Tools        types.List     `tfsdk:"tool_sources"`
	Integrations types.List     `tfsdk:"integrations"`
	Variables    types.Map      `tfsdk:"environment_variables"`
}

type StarterModel struct {
	Name    string `tfsdk:"name"`
	Command string `tfsdk:"command"`
}

func AgentSchema() schema.Schema {
	const (
		emptyJson    = ""
		defaultModel = "gpt-4o"
		defaultImage = "ghcr.io/kubiyabot/kubiya-agent:stable"
	)

	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Required
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the agent",
				MarkdownDescription: "The descriptive name of the agent",
			},
			"runner": schema.StringAttribute{
				Required:            true,
				Description:         "The runner of the agent",
				MarkdownDescription: "The runner used by the agent",
			},
			"description": schema.StringAttribute{
				Required:            true,
				Description:         "The description of the agent",
				MarkdownDescription: "A detailed description of the agent",
			},
			"instructions": schema.StringAttribute{
				Required:            true,
				Description:         "The instructions for the agent",
				MarkdownDescription: "Instructions provided to the agent",
			},

			"is_debug_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Indicate if agent will with debug mode",
				MarkdownDescription: "Indicate if agent will with debug mode",
				Default:             booldefault.StaticBool(false),
			},

			// Required with default values
			"image": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The docker image for the agent",
				MarkdownDescription: "The Docker image used for the agent",
				Default:             stringdefault.StaticString(defaultImage),
			},
			"model": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The LLM model that the agent will run",
				Default:             stringdefault.StaticString(defaultModel),
				MarkdownDescription: "The LLM model used by the agent for its operations",
			},

			// Computed
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the agent",
				MarkdownDescription: "The unique identifier of the agent",
			},
			"owner": schema.StringAttribute{
				Computed:            true,
				Description:         "The owner of the agent",
				MarkdownDescription: "The user who created the agent",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation time of the agent",
				MarkdownDescription: "The timestamp when the agent was created",
			},

			// Optional
			"links": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of links associated with the agent",
				MarkdownDescription: "An array of links related to the agent",
			},
			"tasks": schema.ListAttribute{
				Optional: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"prompt":      types.StringType,
						"description": types.StringType,
					},
				},
				Description:         "A list of tasks associated with the agent",
				MarkdownDescription: "An array of tasks related to the agent",
			},
			"tool_sources": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of tools to be consumed with agents",
				MarkdownDescription: "An array of URL's to pull tools from",
			},
			"users": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of users that have access to this agent",
				MarkdownDescription: "An array of users who have access to this agent",
			},
			"groups": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of groups that have access to this agent",
				MarkdownDescription: "An array of groups who have access to this agent",
			},
			"secrets": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of secrets associated with the agent",
				MarkdownDescription: "An array of secrets related to the agent",
			},
			"sources": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of tools to be consumed with agents",
				MarkdownDescription: "An array of URL's to pull tools from",
			},
			"starters": schema.ListAttribute{
				Optional: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":    types.StringType,
						"command": types.StringType,
					},
				},
				Description:         "A list of starters associated with the agent",
				MarkdownDescription: "An array of starters related to the agent",
			},
			"integrations": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of integrations associated with the agent",
				MarkdownDescription: "An array of integrations related to the agent",
			},
			"environment_variables": schema.MapAttribute{
				Computed:            true,
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "A map of environment variables for the agent",
				MarkdownDescription: "A map of key-value pairs representing environment variables for the agent",
			},
			"workflows": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  defaultString(emptyJson),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonNormalizationModifier(),
				},
			},
		},
	}
}

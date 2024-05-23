package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Image string `tfsdk:"image"`
	Model string `tfsdk:"model"`

	// Required
	Name         string `tfsdk:"name"`
	Runner       string `tfsdk:"runner"`
	Description  string `tfsdk:"description"`
	Instructions string `tfsdk:"instructions"`

	// Optional
	Links        []string       `tfsdk:"links"`
	Tasks        []TaskModel    `tfsdk:"tasks"`
	Users        []string       `tfsdk:"users"`
	Groups       []string       `tfsdk:"groups"`
	Secrets      []string       `tfsdk:"secrets"`
	Starters     []StarterModel `tfsdk:"starters"`
	Integrations []string       `tfsdk:"integrations"`
	Variables    types.Map      `tfsdk:"environment_variables"`
}

type StarterModel struct {
	Name    string `tfsdk:"name"`
	Command string `tfsdk:"command"`
}

func AgentSchema() schema.Schema {
	const (
		defaultModel = "azure/gpt-4"
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

			// Required with default values
			"image": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The docker image for the agent",
				MarkdownDescription: "The Docker image used for the agent",
				Default:             stringdefault.StaticString(defaultImage),
			},
			"model": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					onOfValidator("model", []string{"azure/gpt-4",
						"azure/gpt-4-32k", "azure/gpt-4-32k",
						"azure/gpt-4-turbo-preview", "azure/gpt3-5-turbo-16k"}),
				},
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
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
			"tasks": schema.ListAttribute{
				Optional: true,
				//Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"prompt":      types.StringType,
						"description": types.StringType,
					},
				},
				Description:         "A list of tasks associated with the agent",
				MarkdownDescription: "An array of tasks related to the agent",
				//Default:             listdefault.StaticValue(types.ListNull(types.ObjectType{})),
			},
			"users": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of users that have access to this agent",
				MarkdownDescription: "An array of users who have access to this agent",
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
			"groups": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of groups that have access to this agent",
				MarkdownDescription: "An array of groups who have access to this agent",
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
			"secrets": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of secrets associated with the agent",
				MarkdownDescription: "An array of secrets related to the agent",
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
			"starters": schema.ListAttribute{
				Optional: true,
				//Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":    types.StringType,
						"command": types.StringType,
					},
				},
				Description:         "A list of starters associated with the agent",
				MarkdownDescription: "An array of starters related to the agent",
				//Default:             listdefault.StaticValue(types.ListNull(types.ObjectType{})),
			},
			"integrations": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of integrations associated with the agent",
				MarkdownDescription: "An array of integrations related to the agent",
				Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
			},
			"environment_variables": schema.MapAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A map of environment variables for the agent",
				MarkdownDescription: "A map of key-value pairs representing environment variables for the agent",
				Default:             mapdefault.StaticValue(types.MapNull(types.StringType)),
			},
		},
	}
}

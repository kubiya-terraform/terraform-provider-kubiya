package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentModel struct {
	Id        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`

	Image types.String `tfsdk:"image"`
	Model types.String `tfsdk:"model"`

	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Instructions types.String `tfsdk:"instructions"`

	Links        types.String `tfsdk:"links"`        // list
	Users        types.String `tfsdk:"users"`        // list
	Groups       types.String `tfsdk:"groups"`       // list
	Runners      types.String `tfsdk:"runners"`      // list
	Secrets      types.String `tfsdk:"secrets"`      // list
	Starters     types.String `tfsdk:"starters"`     // list
	Variables    types.String `tfsdk:"env_vars"`     // map[string]string
	Integrations types.String `tfsdk:"integrations"` // list
	Tasks        types.String `tfsdk:"tasks"`        // list
}

func AgentSchema() schema.Schema {
	const (
		empty        = ""
		defaultModel = "azure/gpt-4-32k"
		defaultImage = "ghcr.io/kubiyabot/kubiya-agent:stable"
	)

	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"email":      schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true},
			"created_by": schema.StringAttribute{Computed: true},

			"name":         schema.StringAttribute{Required: true},
			"runners":      schema.StringAttribute{Required: true},
			"description":  schema.StringAttribute{Required: true},
			"instructions": schema.StringAttribute{Required: true},

			"links":        schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"secrets":      schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"starters":     schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"integrations": schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"users":        schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"groups":       schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"env_vars":     schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
			"model":        schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(defaultModel)},
			"image":        schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(defaultImage)},
			"tasks":        schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString(empty)},
		},
	}
}

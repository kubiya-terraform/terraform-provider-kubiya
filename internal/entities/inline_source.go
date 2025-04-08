package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InlineSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Tools  types.String `tfsdk:"tools"`
	Runner types.String `tfsdk:"runner"`
	Config types.String `tfsdk:"dynamic_config"`
}

func InlineSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Computed
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the tool",
				MarkdownDescription: "The unique identifier of the inline source tool",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the inline source",
				MarkdownDescription: "The descriptive type of the inline source",
			},

			// Required
			"tools": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonNormalizationModifier(),
				},
			},

			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the inline source tool",
				MarkdownDescription: "The descriptive name of the inline source",
			},

			// Optional + Computed
			"runner": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "The runner name",
				MarkdownDescription: "The runner name to add for inline source",
			},
			"dynamic_config": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             defaultString("{}"),
				PlanModifiers:       []planmodifier.String{jsonNormalizationModifier()},
				Description:         "The dynamic configuration of the inline source",
				MarkdownDescription: "A map of key-value pairs representing dynamic configuration for the inline source",
			},
		},
	}
}

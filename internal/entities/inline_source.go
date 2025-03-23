package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InlineTool struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Image       types.String `tfsdk:"image"`
	Content     types.String `tfsdk:"content"`
	Description types.String `tfsdk:"description"`
}

type InlineSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Tools  []InlineTool `tfsdk:"tools"`
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

			// Required
			"tools": schema.ListAttribute{
				Required: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"prompt":      types.StringType,
						"image":       types.StringType,
						"content":     types.StringType,
						"description": types.StringType,
					},
				},
				Description:         "A list of tools for inline source",
				MarkdownDescription: "An array of tools for inline source",
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
				Default:             defaultString(),
				Description:         "The runner name",
				MarkdownDescription: "The runner name to add for inline source",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             defaultString(),
				Description:         "The type of the inline source",
				MarkdownDescription: "The descriptive type of the inline source",
			},
			"dynamic_config": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  defaultString(),
				PlanModifiers: []planmodifier.String{
					jsonNormalizationModifier(),
				},
				Description:         "The dynamic configuration of the inline source",
				MarkdownDescription: "A map of key-value pairs representing dynamic configuration for the inline source",
			},
		},
	}
}

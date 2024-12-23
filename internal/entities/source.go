package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	// Required
	Url types.String `tfsdk:"url"`

	// Computed
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DynamicConfig types.String `tfsdk:"dynamic_config"`
	Runner        types.String `tfsdk:"runner"`
}

func SourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Required
			"url": schema.StringAttribute{
				Required:            true,
				Description:         "url path for source",
				MarkdownDescription: "url path for source",
			},

			// Computed
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the source",
				MarkdownDescription: "The unique identifier of the source",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the source",
				MarkdownDescription: "The descriptive name of the source",
			},
			"dynamic_config": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					jsonStringModifier{},
				},
				Description:         "The dynamic configuration of the source",
				MarkdownDescription: "A map of key-value pairs representing dynamic configuration for the source",
			},
			"runner": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "The runner name",
				MarkdownDescription: "The runner name to add the source",
			},
		},
	}
}

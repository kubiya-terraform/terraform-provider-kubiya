package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	// Required
	Url types.String `tfsdk:"url"`

	// Computed
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DynamicConfig types.Map    `tfsdk:"dynamic_config"`
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
			"dynamic_config": schema.MapAttribute{
				Computed:            true,
				Optional:            true,
				ElementType:         types.DynamicType,
				Description:         "The dynamic configuration of the source",
				MarkdownDescription: "A map of key-value pairs representing dynamic configuration for the source",
			},
		},
	}
}

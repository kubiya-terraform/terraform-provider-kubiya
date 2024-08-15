package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	// Computed
	Id types.String `tfsdk:"id"`

	// Required
	Url  types.String `tfsdk:"url"`
	Name types.String `tfsdk:"name"`
}

func SourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Computed
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the source",
				MarkdownDescription: "The unique identifier of the source",
			},
			// Required
			"url": schema.StringAttribute{
				Required:            true,
				Description:         "url path for source",
				MarkdownDescription: "url path for source",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the source",
				MarkdownDescription: "The descriptive name of the source",
			},
		},
	}
}

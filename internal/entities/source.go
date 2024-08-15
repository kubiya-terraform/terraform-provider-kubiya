package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceModel struct {
	// Required
	Url types.String `tfsdk:"url"`

	// Computed
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ToolsCount     types.Int64  `tfsdk:"tools_count"`
	AgentsCount    types.Int64  `tfsdk:"agents_count"`
	WorkflowsCount types.Int64  `tfsdk:"workflows_count"`
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
			"tools_count": schema.Int64Attribute{
				Computed:    true,
				Description: "numbers of tools connected to source",
			},
			"agents_count": schema.Int64Attribute{
				Computed:    true,
				Description: "numbers of agents connected to source",
			},
			"workflows_count": schema.Int64Attribute{
				Computed:    true,
				Description: "numbers of workflows connected to source",
			},
		},
	}
}

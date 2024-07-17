package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KnowledgeModel struct {
	// Computed
	Id    types.String `tfsdk:"id"`
	Owner types.String `tfsdk:"owner"`

	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Groups      types.List   `tfsdk:"groups"`
	Content     types.String `tfsdk:"content"`
	Description types.String `tfsdk:"description"`

	Labels          types.List `tfsdk:"labels"`
	SupportedAgents types.List `tfsdk:"supported_agents"`
}

func KnowledgeSchema() schema.Schema {
	const (
		defaultType = "knowledge"
	)
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the knowledge",
				MarkdownDescription: "The unique identifier of the knowledge",
			},
			"owner": schema.StringAttribute{
				Computed: true,
			},

			// Required
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the knowledge",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the knowledge",
				Default:     stringdefault.StaticString(defaultType),
			},
			"groups": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "A list of user groups with access associated with the knowledge",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "The content of the knowledge",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "The description of the knowledge",
			},

			// Optional
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of labels associated with the knowledge",
			},
			"supported_agents": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "A list of agents associated with the knowledge",
				MarkdownDescription: "An array of agents related to the knowledge",
			},
		},
	}
}

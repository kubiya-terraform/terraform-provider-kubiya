package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExternalKnowledgeModel struct {
	// Computed - UUID from API
	Id types.String `tfsdk:"id"`

	// Required
	Vendor types.String  `tfsdk:"vendor"`
	Config types.Dynamic `tfsdk:"config"`

	// Computed - Additional fields from API response
	Org             types.String `tfsdk:"org"`
	StartDate       types.String `tfsdk:"start_date"`
	IntegrationType types.String `tfsdk:"integration_type"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func ExternalKnowledgeSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The UUID of the external knowledge integration",
				MarkdownDescription: "The unique identifier of the external knowledge integration returned from the API",
			},
			"vendor": schema.StringAttribute{
				Required:    true,
				Description: "The vendor/provider for the knowledge integration (e.g., 'slack', 'confluence', 'notion')",
			},
			"config": schema.DynamicAttribute{
				Required:    true,
				Description: "Dynamic configuration for the vendor. Supports maps with strings, lists, and other types. Examples: For Slack: {'channel_ids': ['C1234567890', 'C0987654321']}. For Confluence: {'space_key': 'DEV', 'page_id': '123456'}",
			},
			"org": schema.StringAttribute{
				Computed:    true,
				Description: "The organization associated with the integration",
			},
			"start_date": schema.StringAttribute{
				Computed:    true,
				Description: "The start date of the integration",
			},
			"integration_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of integration (matches the vendor field)",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the integration was created",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the integration was last updated",
			},
		},
	}
}

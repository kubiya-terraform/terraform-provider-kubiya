package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebhookModel struct {
	Id          types.String `tfsdk:"id"`
	Url         types.String `tfsdk:"url"`
	Name        types.String `tfsdk:"name"`
	Agent       types.String `tfsdk:"agent"`
	Filter      types.String `tfsdk:"filter"`
	Source      types.String `tfsdk:"source"`
	Prompt      types.String `tfsdk:"prompt"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	Destination types.String `tfsdk:"destination"`
}

func WebhookSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"created_at":  schema.StringAttribute{Computed: true},
			"created_by":  schema.StringAttribute{Required: true},
			"url":         schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"agent":       schema.StringAttribute{Required: true},
			"source":      schema.StringAttribute{Required: true},
			"prompt":      schema.StringAttribute{Required: true},
			"destination": schema.StringAttribute{Required: true},
			"filter":      schema.StringAttribute{Optional: true},
		},
	}
}

package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecretModel struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
}

func SecretSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":        schema.StringAttribute{Required: true},
			"value":       schema.StringAttribute{Required: true, Sensitive: true},
			"description": schema.StringAttribute{Optional: true},
			"created_at":  schema.StringAttribute{Computed: true},
			"created_by":  schema.StringAttribute{Computed: true},
		},
	}
}

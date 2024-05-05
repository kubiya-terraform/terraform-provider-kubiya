package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderConfig struct {
	UserKey types.String `tfsdk:"user_key"`
}

func ProviderConfigSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_key": schema.StringAttribute{
				Required:  true,
				Sensitive: false,
			},
		},
	}
}

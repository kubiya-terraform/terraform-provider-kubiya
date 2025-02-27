package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunnerModel struct {
	Name       types.String `tfsdk:"name"`
	RunnerType types.String `tfsdk:"runner_type"`
}

func RunnerSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":        schema.StringAttribute{Required: true},
			"runner_type": schema.StringAttribute{Computed: true},
		},
	}
}

package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunnerModel struct {
	Key     types.String `tfsdk:"key"`
	Url     types.String `tfsdk:"url"`
	Name    types.String `tfsdk:"name"`
	Path    types.String `tfsdk:"path"`
	Subject types.String `tfsdk:"subject"`
}

func RunnerSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":    schema.StringAttribute{Required: true},
			"key":     schema.StringAttribute{Computed: true},
			"url":     schema.StringAttribute{Computed: true},
			"subject": schema.StringAttribute{Computed: true},
			"path":    schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationModel struct {
	ID          types.String  `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Configs     []ConfigModel `tfsdk:"configs"`
	AuthType    types.String  `tfsdk:"auth_type"`
	Description types.String  `tfsdk:"description"`
	Type        types.String  `tfsdk:"integration_type"`
}

type ConfigModel struct {
	Name           types.String `tfsdk:"name"`
	IsDefault      types.Bool   `tfsdk:"is_default"`
	VendorSpecific types.Map    `tfsdk:"vendor_specific"`
}

func IntegrationSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the integration",
				MarkdownDescription: "The unique identifier of the integration",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the integration",
				MarkdownDescription: "The name of the integration",
			},
			"integration_type": schema.StringAttribute{
				Required:            true,
				Description:         "The type of the integration",
				MarkdownDescription: "The type of the integration (e.g., aws, aws_organization, gcp, azure, jira, confluence)",
			},
			"auth_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Description:         "The authentication type of the integration",
				MarkdownDescription: "The authentication type of the integration (e.g., per_user, global)",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "The description of the integration",
				MarkdownDescription: "A description of the integration",
			},
			"configs": schema.ListAttribute{
				Required: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":       types.StringType,
						"is_default": types.BoolType,
						"vendor_specific": types.MapType{
							ElemType: types.StringType,
						},
					},
				},
			},
		},
	}
}

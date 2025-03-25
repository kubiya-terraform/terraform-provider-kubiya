package entities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	TeamName    types.String `tfsdk:"team_name"`
	Method      types.String `tfsdk:"method"`
}

// Custom validator for team_name based on method
type teamNameValidator struct{}

func (v teamNameValidator) Description(ctx context.Context) string {
	return "Validates team_name is provided when method is 'teams'"
}

func (v teamNameValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that `team_name` is provided when `method` is set to `teams`"
}

func (v teamNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if team_name is not being configured
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	methodPath := path.Root("method")
	var method types.String
	diags := req.Config.GetAttribute(ctx, methodPath, &method)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If method is teams, team_name must be provided and not empty
	if !method.IsNull() && !method.IsUnknown() && method.ValueString() == "teams" {
		if req.ConfigValue.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing TeamName",
				"When Method is 'teams', TeamName is required and cannot be empty",
			)
		}
	}
}

func WebhookSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"created_at":  schema.StringAttribute{Computed: true},
			"created_by":  schema.StringAttribute{Computed: true},
			"url":         schema.StringAttribute{Computed: true},
			"name":        schema.StringAttribute{Required: true},
			"agent":       schema.StringAttribute{Required: true},
			"source":      schema.StringAttribute{Required: true},
			"prompt":      schema.StringAttribute{Required: true},
			"destination": schema.StringAttribute{Required: true},
			"filter":      schema.StringAttribute{Optional: true},
			"method": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("Slack"),
			},
			"team_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					teamNameValidator{},
				},
			},
		},
	}
}

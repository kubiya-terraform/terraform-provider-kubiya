package entities

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

	Runner   types.String `tfsdk:"runner"`
	Workflow types.String `tfsdk:"workflow"`
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

	// Only validate team_name if method is "teams"
	if !method.IsNull() && !method.IsUnknown() && strings.EqualFold(method.ValueString(), "teams") {
		if req.ConfigValue.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing TeamName",
				"When Method is 'teams', TeamName is required and cannot be empty",
			)
		}
	}
}

// Add a new validator for destination
type destinationValidator struct{}

func (v destinationValidator) Description(ctx context.Context) string {
	return "Validates destination is provided when method is not 'http'"
}

func (v destinationValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that `destination` is provided when `method` is not `http`"
}

func (v destinationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	methodPath := path.Root("method")
	var method types.String
	diags := req.Config.GetAttribute(ctx, methodPath, &method)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If method is not "http", destination is required
	if !method.IsNull() && !method.IsUnknown() && !strings.EqualFold(method.ValueString(), "http") {
		if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Destination",
				"Destination is required when Method is not 'http'",
			)
		}
	}
}

func WebhookSchema() schema.Schema {
	const emptyJson = ""
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true},
			"created_by": schema.StringAttribute{Computed: true},
			"url":        schema.StringAttribute{Computed: true},
			"name":       schema.StringAttribute{Required: true},
			"agent":      schema.StringAttribute{Optional: true},
			"runner":     schema.StringAttribute{Optional: true},
			"workflow": schema.StringAttribute{Computed: true,
				Optional: true,
				Default:  defaultString(emptyJson),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonNormalizationModifier(),
				}},
			"source": schema.StringAttribute{Optional: true},
			"prompt": schema.StringAttribute{Required: true},
			"destination": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					destinationValidator{},
				},
			},
			"filter": schema.StringAttribute{Optional: true},
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

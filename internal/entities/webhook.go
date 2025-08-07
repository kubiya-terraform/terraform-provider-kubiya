package entities

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WebhookModel represents the Terraform resource model for a webhook.
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
	Runner      types.String `tfsdk:"runner"`
	Workflow    types.String `tfsdk:"workflow"`
}

// jsonValidator ensures the provided string is valid JSON.
type jsonValidator struct{}

func (v jsonValidator) Description(_ context.Context) string {
	return "Ensures the provided string is valid JSON"
}

func (v jsonValidator) MarkdownDescription(_ context.Context) string {
	return "Ensures that the `workflow` field contains a valid JSON string"
}

func (v jsonValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	var jsonRaw json.RawMessage
	if err := json.Unmarshal([]byte(value), &jsonRaw); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON",
			"The provided value is not a valid JSON string: "+err.Error(),
		)
	}
}

// teamNameValidator validates team_name based on method.
type teamNameValidator struct{}

func (v teamNameValidator) Description(_ context.Context) string {
	return "Validates team_name is provided when method is 'teams'"
}

func (v teamNameValidator) MarkdownDescription(_ context.Context) string {
	return "Validates that `team_name` is provided when `method` is set to `teams`"
}

func (v teamNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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

// destinationValidator validates destination based on method.
type destinationValidator struct{}

func (v destinationValidator) Description(_ context.Context) string {
	return "Validates destination is provided when method is not 'http'"
}

func (v destinationValidator) MarkdownDescription(_ context.Context) string {
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

// jsonNormalizationModifier ensures consistent JSON string normalization.
type jsonNormalizationModifierWH struct{}

func (m jsonNormalizationModifierWH) Description(_ context.Context) string {
	return "Normalizes JSON strings to ensure consistent formatting"
}

func (m jsonNormalizationModifierWH) MarkdownDescription(_ context.Context) string {
	return "Ensures that the JSON string in the `workflow` field is consistently formatted by parsing and re-encoding it"
}

func (m jsonNormalizationModifierWH) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	value := req.PlanValue.ValueString()
	if value == "" {
		return
	}

	// Parse and normalize the JSON string
	var jsonRaw json.RawMessage
	if err := json.Unmarshal([]byte(value), &jsonRaw); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid JSON in Plan",
			"Failed to parse JSON string: "+err.Error(),
		)
		return
	}

	// Re-encode with consistent formatting (no indentation, standard escaping)
	normalized, err := json.Marshal(jsonRaw)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"JSON Normalization Failed",
			"Failed to normalize JSON string: "+err.Error(),
		)
		return
	}

	resp.PlanValue = types.StringValue(string(normalized))
}

// nullIfEmptyModifier ensures empty strings are treated as null.
type nullIfEmptyModifier struct{}

func (m nullIfEmptyModifier) Description(_ context.Context) string {
	return "Converts empty strings to null"
}

func (m nullIfEmptyModifier) MarkdownDescription(_ context.Context) string {
	return "Ensures that empty strings are treated as null in the plan and state"
}

func (m nullIfEmptyModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	if req.PlanValue.ValueString() == "" {
		resp.PlanValue = types.StringNull()
	}
}

// defaultStaticString creates a default string plan modifier.
func defaultStaticString(value string) defaults.String {
	return stringdefault.StaticString(value)
}

// WebhookSchema defines the schema for the webhook resource.
func WebhookSchema() schema.Schema {
	const emptyJson = ""
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the webhook",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the webhook was created",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User or entity that created the webhook",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "URL for the webhook endpoint",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the webhook",
			},
			"agent": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     defaultStaticString(""),
				Description: "Agent associated with the webhook",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					nullIfEmptyModifier{},
				},
			},
			"runner": schema.StringAttribute{
				Optional:    true,
				Description: "Runner configuration for the webhook",
			},
			"workflow": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Default:     defaultStaticString(emptyJson),
				Description: "JSON string defining the workflow configuration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					jsonNormalizationModifierWH{},
				},
				Validators: []validator.String{
					jsonValidator{},
				},
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "Source of the webhook trigger",
			},
			"prompt": schema.StringAttribute{
				Required:    true,
				Description: "Prompt or trigger condition for the webhook",
			},
			"destination": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     defaultStaticString("webhook"),
				Description: "Destination for the webhook payload",
				Validators: []validator.String{
					destinationValidator{},
				},
			},
			"filter": schema.StringAttribute{
				Optional:    true,
				Description: "Filter criteria for the webhook",
			},
			"method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     defaultStaticString("Slack"),
				Description: "Method for delivering the webhook (e.g., Slack, teams, http)",
			},
			"team_name": schema.StringAttribute{
				Optional:    true,
				Description: "Team name for the webhook when method is 'teams'",
				Validators: []validator.String{
					teamNameValidator{},
				},
			},
		},
	}
}

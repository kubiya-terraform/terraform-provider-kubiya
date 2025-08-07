package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*webhookResource)(nil)
	_ resource.ResourceWithConfigure = (*webhookResource)(nil)
)

type webhookResource struct {
	name   string
	client *clients.Client
}

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	if err := r.client.ReadWebhook(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"webhook not found",
			fmt.Sprintf("webhook by name: %s not found. Error: ", state.Name)+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.WebhookModel
	var state entities.WebhookModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	if !plan.Id.IsUnknown() && !plan.Id.IsNull() {
		updatedState.Id = plan.Id
	}
	if !plan.Url.IsUnknown() && !plan.Url.IsNull() {
		updatedState.Url = plan.Url
	}
	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		updatedState.Name = plan.Name
	}
	if !plan.Agent.IsUnknown() && !plan.Agent.IsNull() {
		updatedState.Agent = plan.Agent
	}
	if !plan.Prompt.IsUnknown() && !plan.Prompt.IsNull() {
		updatedState.Prompt = plan.Prompt
	}
	if !plan.Source.IsUnknown() && !plan.Source.IsNull() {
		updatedState.Source = plan.Source
	}
	if !plan.Filter.IsUnknown() && !plan.Filter.IsNull() {
		updatedState.Filter = plan.Filter
	}
	if !plan.Destination.IsUnknown() && !plan.Destination.IsNull() {
		updatedState.Destination = plan.Destination
	}

	if err := r.client.UpdateWebhook(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			"failed to update webhook",
			"failed to update webhook. Error: "+err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.WebhookSchema()
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Normalize workflow JSON before sending to backend
	workflow := plan.Workflow.ValueString()
	if workflow != "" {
		var jsonRaw json.RawMessage
		if err := json.Unmarshal([]byte(workflow), &jsonRaw); err != nil {
			resp.Diagnostics.AddError(
				"Invalid JSON in Workflow",
				"Failed to parse workflow JSON: "+err.Error(),
			)
			return
		}
		normalized, err := json.Marshal(jsonRaw)
		if err != nil {
			resp.Diagnostics.AddError(
				"JSON Normalization Failed",
				"Failed to normalize workflow JSON: "+err.Error(),
			)
			return
		}
		plan.Workflow = types.StringValue(string(normalized))
	}

	// Handle agent field: ensure empty string is treated as null
	if plan.Agent.ValueString() == "" {
		plan.Agent = types.StringNull()
	}

	// Call backend API to create webhook
	state, err := r.client.CreateWebhook(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Webhook",
			"Failed to create webhook. Error: "+err.Error(),
		)
		return
	}

	// Normalize workflow JSON from backend response
	if state.Workflow.ValueString() != "" {
		var jsonRaw json.RawMessage
		if err := json.Unmarshal([]byte(state.Workflow.ValueString()), &jsonRaw); err == nil {
			normalized, err := json.Marshal(jsonRaw)
			if err == nil {
				state.Workflow = types.StringValue(string(normalized))
			} else {
				resp.Diagnostics.AddWarning(
					"JSON Normalization Warning",
					"Failed to normalize workflow JSON from backend: "+err.Error(),
				)
			}
		} else {
			resp.Diagnostics.AddWarning(
				"Invalid JSON in Backend Response",
				"Backend returned invalid workflow JSON: "+err.Error(),
			)
		}
	}

	// Handle agent field from backend response: convert empty string to null
	if state.Agent.ValueString() == "" {
		state.Agent = types.StringNull()
	}

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	if err := r.client.DeleteWebhook(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"failed to delete webhook",
			"failed to delete webhook. Error: "+err.Error(),
		)
	}
}

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "webhook"
		r.client = client
	}
}

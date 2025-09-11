package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/getsentry/sentry-go"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"

	kubiyasentry "terraform-provider-kubiya/internal/sentry"
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
	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Reading webhook resource", sentry.LevelInfo, nil)

	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Update span with resource ID if available
	id := ""
	if !state.Id.IsNull() {
		id = state.Id.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
		kubiyasentry.SetSpanData(span, "webhook.id", id)
	}

	// Read API call logic
	if err := r.client.ReadWebhook(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_webhook", id, "read", err)

		resp.Diagnostics.AddError(
			"webhook not found",
			fmt.Sprintf("webhook by name: %s not found. Error: ", state.Name)+err.Error(),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Updating webhook resource", sentry.LevelInfo, nil)

	var plan entities.WebhookModel
	var state entities.WebhookModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	id := state.Id.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
	kubiyasentry.SetSpanData(span, "webhook.id", id)

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
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_webhook", id, "update", err)

		resp.Diagnostics.AddError(
			"failed to update webhook",
			"failed to update webhook. Error: "+err.Error(),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.WebhookSchema()
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Creating webhook resource", sentry.LevelInfo, nil)

	var plan entities.WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Add plan data to span
	if !plan.Name.IsNull() {
		kubiyasentry.SetSpanData(span, "webhook.name", plan.Name.ValueString())
	}

	// Normalize workflow JSON before sending to backend
	workflow := plan.Workflow.ValueString()
	if workflow != "" {
		var jsonRaw json.RawMessage
		if err := json.Unmarshal([]byte(workflow), &jsonRaw); err != nil {
			kubiyasentry.RecordError(ctx, err)
			CaptureResourceError(ctx, "kubiya_webhook", "", "create", err)

			resp.Diagnostics.AddError(
				"Invalid JSON in Workflow",
				"Failed to parse workflow JSON: "+err.Error(),
			)
			return
		}
		normalized, err := json.Marshal(jsonRaw)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			CaptureResourceError(ctx, "kubiya_webhook", "", "create", err)

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
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_webhook", "", "create", err)

		resp.Diagnostics.AddError(
			"Failed to Create Webhook",
			"Failed to create webhook. Error: "+err.Error(),
		)
		return
	}

	// Update span with created resource ID
	if state != nil && !state.Id.IsNull() {
		id := state.Id.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
		kubiyasentry.SetSpanData(span, "webhook.created_id", id)
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

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Deleting webhook resource", sentry.LevelInfo, nil)

	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Update span with resource ID
	id := state.Id.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
	kubiyasentry.SetSpanData(span, "webhook.id", id)

	// Delete API call logic
	if err := r.client.DeleteWebhook(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_webhook", id, "delete", err)

		resp.Diagnostics.AddError(
			"failed to delete webhook",
			"failed to delete webhook. Error: "+err.Error(),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	kubiyasentry.AddBreadcrumb("resource", "Successfully deleted webhook "+id, sentry.LevelInfo, nil)
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

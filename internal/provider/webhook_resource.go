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
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for webhook read", "error", "diagnostics error")
		return
	}

	// Update span with resource ID if available
	id := ""
	name := ""
	if !state.Id.IsNull() {
		id = state.Id.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
		kubiyasentry.SetSpanData(span, "webhook.id", id)
	}
	if !state.Name.IsNull() {
		name = state.Name.ValueString()
	}

	// Log read operation
	logger.Debug("Reading webhook resource", "webhook_id", id, "webhook_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading webhook resource", sentry.LevelDebug, map[string]interface{}{"webhook_id": id, "webhook_name": name})

	// Read API call logic
	if err := r.client.ReadWebhook(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read webhook", "webhook_id", id, "webhook_name", name, "error", err)
		resp.Diagnostics.AddError(
			"webhook not found",
			fmt.Sprintf("webhook by name: %s not found. Error: ", state.Name)+err.Error(),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read webhook", "webhook_id", id, "webhook_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.WebhookModel
	var state entities.WebhookModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan or state for webhook update", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
	kubiyasentry.SetSpanData(span, "webhook.id", id)

	// Log update operation
	logger.Info("Updating webhook resource", "webhook_id", id, "webhook_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating webhook resource", sentry.LevelInfo, map[string]interface{}{"webhook_id": id, "webhook_name": name})

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
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update webhook", "webhook_id", id, "webhook_name", name, "error", err)
		resp.Diagnostics.AddError(
			"failed to update webhook",
			"failed to update webhook. Error: "+err.Error(),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated webhook", "webhook_id", id, "webhook_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.WebhookSchema()
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for webhook create", "error", "diagnostics error")
		return
	}

	// Add plan data to span
	name := ""
	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
		kubiyasentry.SetSpanData(span, "webhook.name", name)
	}

	// Log create operation
	logger.Info("Creating webhook resource", "webhook_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating webhook resource", sentry.LevelInfo, map[string]interface{}{"webhook_name": name})

	// Normalize workflow JSON before sending to backend
	workflow := plan.Workflow.ValueString()
	if workflow != "" {
		var jsonRaw json.RawMessage
		if err := json.Unmarshal([]byte(workflow), &jsonRaw); err != nil {
			kubiyasentry.RecordError(ctx, err)
			kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
			logger.Error("Invalid JSON in workflow", "webhook_name", name, "error", err)
			resp.Diagnostics.AddError(
				"Invalid JSON in Workflow",
				"Failed to parse workflow JSON: "+err.Error(),
			)
			return
		}
		normalized, err := json.Marshal(jsonRaw)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
			logger.Error("Failed to normalize workflow JSON", "webhook_name", name, "error", err)
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
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create webhook", "webhook_name", name, "error", err)
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
				logger.Warn("Failed to normalize workflow JSON from backend", "error", err)
				resp.Diagnostics.AddWarning(
					"JSON Normalization Warning",
					"Failed to normalize workflow JSON from backend: "+err.Error(),
				)
			}
		} else {
			logger.Warn("Backend returned invalid workflow JSON", "error", err)
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

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created webhook", "webhook_name", name, "webhook_id", state.Id.ValueString())
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_webhook", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for webhook delete", "error", "diagnostics error")
		return
	}

	// Update span with resource ID
	id := state.Id.ValueString()
	name := state.Name.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, id)
	kubiyasentry.SetSpanData(span, "webhook.id", id)

	// Log delete operation
	logger.Info("Deleting webhook resource", "webhook_id", id, "webhook_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting webhook resource", sentry.LevelInfo, map[string]interface{}{"webhook_id": id, "webhook_name": name})

	// Delete API call logic
	if err := r.client.DeleteWebhook(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete webhook", "webhook_id", id, "webhook_name", name, "error", err)
		resp.Diagnostics.AddError(
			"failed to delete webhook",
			"failed to delete webhook. Error: "+err.Error(),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted webhook", "webhook_id", id, "webhook_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Successfully deleted webhook "+id, sentry.LevelInfo, nil)
}

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			logger.Error("Failed to configure webhook resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "webhook"
		r.client = client
		logger.Debug("Successfully configured webhook resource")
	}
}

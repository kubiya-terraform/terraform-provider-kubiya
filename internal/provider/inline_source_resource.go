package provider

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

var (
	_ resource.Resource                = (*inlineSourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*inlineSourceResource)(nil)
	_ resource.ResourceWithImportState = (*inlineSourceResource)(nil)
)

type inlineSourceResource struct {
	name   string
	client *clients.Client
}

func NewInlineSourceResource() resource.Resource {
	return &inlineSourceResource{}
}

func (r *inlineSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_inline_source", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.InlineSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for inline_source read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log read operation
	logger.Debug("Reading inline_source resource", "inline_source_id", id, "inline_source_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading inline_source resource", sentry.LevelDebug, map[string]interface{}{"inline_source_id": id, "inline_source_name": name})

	updatedState, err := r.client.ReadInlineSource(ctx, id)
	if err != nil || updatedState == nil {
		if err == nil {
			err = fmt.Errorf("inline_source %s not found", id)
		}
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read inline_source", "inline_source_id", id, "inline_source_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read inline_source", "inline_source_id", id, "inline_source_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *inlineSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.InlineSourceSchema()
}

func (r *inlineSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_inline_source", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.InlineSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for inline_source create", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()
	tools := plan.Tools.ValueString()
	workflow := plan.Workflows.ValueString()

	// Log create operation
	logger.Info("Creating inline_source resource", "inline_source_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating inline_source resource", sentry.LevelInfo, map[string]interface{}{"inline_source_name": name})

	if tools == "{}" && workflow == "{}" {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Invalid inline_source configuration", "inline_source_name", name, "error", "tools and workflows cannot be empty")
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, "tools and workflows cannot be empty"),
		)
		return
	}

	state, err := r.client.CreateInlineSource(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create inline_source", "inline_source_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created inline_source", "inline_source_name", name, "inline_source_id", state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *inlineSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_inline_source", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.InlineSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for inline_source delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log delete operation
	logger.Info("Deleting inline_source resource", "inline_source_id", id, "inline_source_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting inline_source resource", sentry.LevelInfo, map[string]interface{}{"inline_source_id": id, "inline_source_name": name})

	if err := r.client.DeleteInlineSource(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete inline_source", "inline_source_id", id, "inline_source_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted inline_source", "inline_source_id", id, "inline_source_name", name)
}

func (r *inlineSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_inline_source", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.InlineSourceModel
	var state entities.InlineSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for inline_source update", "error", "diagnostics error")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for inline_source update", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log update operation
	logger.Info("Updating inline_source resource", "inline_source_id", id, "inline_source_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating inline_source resource", sentry.LevelInfo, map[string]interface{}{"inline_source_id": id, "inline_source_name": name})

	// Destroy the current resource
	if err := r.client.DeleteInlineSource(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete inline_source during update", "inline_source_id", id, "inline_source_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	// Re-create the resource using the plan
	newState, err := r.client.CreateInlineSource(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create inline_source during update", "inline_source_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated inline_source", "inline_source_name", name, "inline_source_id", newState.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *inlineSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inline_source"
}

func (r *inlineSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure inline_source resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "inline_source"
		r.client = client
		logger.Debug("Successfully configured inline_source resource")
	}
}

func (r *inlineSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

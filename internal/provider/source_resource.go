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
	_ resource.Resource                = (*sourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*sourceResource)(nil)
	_ resource.ResourceWithImportState = (*sourceResource)(nil)
)

type sourceResource struct {
	name   string
	client *clients.Client
}

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

func (r *sourceResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_source", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for source read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log read operation
	logger.Debug("Reading source resource", "source_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Reading source resource", sentry.LevelDebug, map[string]interface{}{"source_id": id})

	updatedState, err := r.client.ReadSource(ctx, id)
	if err != nil || updatedState == nil {
		if err == nil {
			err = fmt.Errorf("source %s not found", id)
		}
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read source", "source_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read source", "source_id", id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *sourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SourceSchema()
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_source", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for source delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log delete operation
	logger.Info("Deleting source resource", "source_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Deleting source resource", sentry.LevelInfo, map[string]interface{}{"source_id": id})

	if err := r.client.DeleteSource(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete source", "source_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted source", "source_id", id)
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_source", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.SourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for source create", "error", "diagnostics error")
		return
	}

	// Log create operation
	logger.Info("Creating source resource")
	kubiyasentry.AddBreadcrumb("resource", "Creating source resource", sentry.LevelInfo, nil)

	state, err := r.client.CreateSource(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create source", "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created source", "source_id", state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure source resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "source"
		r.client = client
		logger.Debug("Successfully configured source resource")
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

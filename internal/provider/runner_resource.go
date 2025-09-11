package provider

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

var (
	_ resource.Resource              = (*runnerResource)(nil)
	_ resource.ResourceWithConfigure = (*runnerResource)(nil)
)

type runnerResource struct {
	name   string
	client *clients.Client
}

func NewRunnerResource() resource.Resource {
	return &runnerResource{}
}

func (r *runnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.RunnerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for runner read", "error", "diagnostics error")
		return
	}

	// Update span with resource name if available
	name := ""
	if !state.Name.IsNull() {
		name = state.Name.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, name)
		kubiyasentry.SetSpanData(span, "runner.name", name)
	}

	// Log read operation
	logger.Debug("Reading runner resource", "runner_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading runner resource", sentry.LevelDebug, map[string]interface{}{"runner_name": name})

	if err := r.client.ReadRunner(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read runner", "runner_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read runner", "runner_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *runnerResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *runnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.RunnerSchema()
}

func (r *runnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.RunnerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for runner create", "error", "diagnostics error")
		return
	}

	// Add plan data to span
	name := ""
	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
		kubiyasentry.SetSpanData(span, "runner.name", name)
	}

	// Log create operation
	logger.Info("Creating runner resource", "runner_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating runner resource", sentry.LevelInfo, map[string]interface{}{"runner_name": name})

	state, err := r.client.CreateRunner(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create runner", "runner_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	// Update span with created resource name
	if state != nil && !state.Name.IsNull() {
		createdName := state.Name.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, createdName)
		kubiyasentry.SetSpanData(span, "runner.created_name", createdName)
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created runner", "runner_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *runnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.RunnerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for runner delete", "error", "diagnostics error")
		return
	}

	// Update span with resource name
	name := state.Name.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, name)
	kubiyasentry.SetSpanData(span, "runner.name", name)

	// Log delete operation
	logger.Info("Deleting runner resource", "runner_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting runner resource", sentry.LevelInfo, map[string]interface{}{"runner_name": name})

	if err := r.client.DeleteRunner(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete runner", "runner_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted runner", "runner_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Successfully deleted runner "+name, sentry.LevelInfo, nil)
}

func (r *runnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner"
}

func (r *runnerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure runner resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "runner"
		r.client = client
		logger.Debug("Successfully configured runner resource")
	}
}

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
	_ resource.Resource                = (*scheduledTaskResource)(nil)
	_ resource.ResourceWithConfigure   = (*scheduledTaskResource)(nil)
	_ resource.ResourceWithImportState = (*scheduledTaskResource)(nil)
)

type scheduledTaskResource struct {
	name   string
	client *clients.Client
}

func NewScheduledTaskResource() resource.Resource {
	return &scheduledTaskResource{}
}

func (r *scheduledTaskResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *scheduledTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_scheduled_task", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.ScheduledTaskModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for scheduled_task read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log read operation
	logger.Debug("Reading scheduled_task resource", "scheduled_task_id", id, "scheduled_task_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Reading scheduled_task resource", sentry.LevelDebug, map[string]interface{}{"scheduled_task_id": id})

	updatedState, err := r.client.ReadScheduledTask(ctx, id)
	if err != nil || updatedState == nil {
		if err == nil {
			err = fmt.Errorf("scheduled_task %s not found", id)
		}
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read scheduled_task", "scheduled_task_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read scheduled_task", "scheduled_task_id", id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *scheduledTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ScheduledTaskSchema()
}

func (r *scheduledTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_scheduled_task", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.ScheduledTaskModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for scheduled_task delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log delete operation
	logger.Info("Deleting scheduled_task resource", "scheduled_task_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Deleting scheduled_task resource", sentry.LevelInfo, map[string]interface{}{"scheduled_task_id": id})

	if err := r.client.DeleteScheduledTask(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete scheduled_task", "scheduled_task_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted scheduled_task", "scheduled_task_id", id)
}

func (r *scheduledTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_scheduled_task", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.ScheduledTaskModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for scheduled_task create", "error", "diagnostics error")
		return
	}

	// Log create operation
	logger.Info("Creating scheduled_task resource")
	kubiyasentry.AddBreadcrumb("resource", "Creating scheduled_task resource", sentry.LevelInfo, map[string]interface{}{})

	state, err := r.client.CreateScheduledTask(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create scheduled_task", "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created scheduled_task", "scheduled_task_id", state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scheduledTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduled_task"
}

func (r *scheduledTaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure scheduled_task resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "scheduled_task"
		r.client = client
		logger.Debug("Successfully configured scheduled_task resource")
	}
}

func (r *scheduledTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

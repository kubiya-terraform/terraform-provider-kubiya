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
	_ resource.Resource                = (*integrationResource)(nil)
	_ resource.ResourceWithConfigure   = (*integrationResource)(nil)
	_ resource.ResourceWithImportState = (*integrationResource)(nil)
)

type integrationResource struct {
	name   string
	client *clients.Client
}

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_integration", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.IntegrationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for integration read", "error", "diagnostics error")
		return
	}

	id := state.ID.ValueString()
	name := state.Name.ValueString()

	// Log read operation
	logger.Debug("Reading integration resource", "integration_id", id, "integration_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading integration resource", sentry.LevelDebug, map[string]interface{}{"integration_id": id, "integration_name": name})

	updatedState, err := r.client.ReadIntegration(ctx, id)
	if err != nil || updatedState == nil {
		if err == nil {
			err = fmt.Errorf("integration %s not found", id)
		}
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read integration", "integration_id", id, "integration_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read integration", "integration_id", id, "integration_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.IntegrationSchema()
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_integration", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.IntegrationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for integration delete", "error", "diagnostics error")
		return
	}

	id := state.ID.ValueString()
	name := state.Name.ValueString()

	// Log delete operation
	logger.Info("Deleting integration resource", "integration_id", id, "integration_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting integration resource", sentry.LevelInfo, map[string]interface{}{"integration_id": id, "integration_name": name})

	if err := r.client.DeleteIntegration(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete integration", "integration_id", id, "integration_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted integration", "integration_id", id, "integration_name", name)
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_integration", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.IntegrationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for integration create", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()

	// Log create operation
	logger.Info("Creating integration resource", "integration_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating integration resource", sentry.LevelInfo, map[string]interface{}{"integration_name": name})

	state, err := r.client.CreateIntegration(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create integration", "integration_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created integration", "integration_name", name, "integration_id", state.ID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_integration", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.IntegrationModel
	var state entities.IntegrationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan or state for integration update", "error", "diagnostics error")
		return
	}

	id := state.ID.ValueString()
	name := state.Name.ValueString()

	// Log update operation
	logger.Info("Updating integration resource", "integration_id", id, "integration_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating integration resource", sentry.LevelInfo, map[string]interface{}{"integration_id": id, "integration_name": name})

	updatedState := state
	updatedState.Configs = plan.Configs

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatedState.Name = plan.Name
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		updatedState.Type = plan.Type
	}

	if !plan.AuthType.IsNull() && !plan.AuthType.IsUnknown() {
		updatedState.AuthType = plan.AuthType
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatedState.Description = plan.Description
	}

	if err := r.client.UpdateIntegration(ctx, &updatedState); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update integration", "integration_id", id, "integration_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated integration", "integration_id", id, "integration_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure integration resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "integration"
		r.client = client
		logger.Debug("Successfully configured integration resource")
	}
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

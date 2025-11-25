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
	_ resource.Resource              = (*externalKnowledgeResource)(nil)
	_ resource.ResourceWithConfigure = (*externalKnowledgeResource)(nil)
)

type externalKnowledgeResource struct {
	name   string
	client *clients.Client
}

func NewExternalKnowledgeResource() resource.Resource {
	return &externalKnowledgeResource{}
}

func (r *externalKnowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.ExternalKnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for external_knowledge read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log read operation
	logger.Debug("Reading external_knowledge resource", "external_knowledge_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Reading external_knowledge resource", sentry.LevelDebug, map[string]interface{}{"external_knowledge_id": id})

	if err := r.client.ReadExternalKnowledge(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read external_knowledge", "external_knowledge_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read external_knowledge", "external_knowledge_id", id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *externalKnowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ExternalKnowledgeSchema()
}

func (r *externalKnowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.ExternalKnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for external_knowledge delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log delete operation
	logger.Info("Deleting external_knowledge resource", "external_knowledge_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Deleting external_knowledge resource", sentry.LevelInfo, map[string]interface{}{"external_knowledge_id": id})

	if err := r.client.DeleteExternalKnowledge(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete external_knowledge", "external_knowledge_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted external_knowledge", "external_knowledge_id", id)
}

func (r *externalKnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.ExternalKnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for external_knowledge create", "error", "diagnostics error")
		return
	}

	// Log create operation
	logger.Info("Creating external_knowledge resource")
	kubiyasentry.AddBreadcrumb("resource", "Creating external_knowledge resource", sentry.LevelInfo, map[string]interface{}{})

	state, err := r.client.CreateExternalKnowledge(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create external_knowledge", "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created external_knowledge", "external_knowledge_id", state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *externalKnowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.ExternalKnowledgeModel
	var state entities.ExternalKnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan or state for external_knowledge update", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log update operation
	logger.Info("Updating external_knowledge resource", "external_knowledge_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Updating external_knowledge resource", sentry.LevelInfo, map[string]interface{}{"external_knowledge_id": id})

	updatedState := state

	// Update vendor if it has changed
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		updatedState.Vendor = plan.Vendor
	}

	// Update config if it has changed
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		updatedState.Config = plan.Config
	}

	if err := r.client.UpdateExternalKnowledge(ctx, &updatedState); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update external_knowledge", "external_knowledge_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated external_knowledge", "external_knowledge_id", id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *externalKnowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_knowledge"
}

func (r *externalKnowledgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure external_knowledge resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "external_knowledge"
		r.client = client
		logger.Debug("Successfully configured external_knowledge resource")
	}
}

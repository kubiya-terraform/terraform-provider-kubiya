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
	_ resource.Resource              = (*knowledgeResource)(nil)
	_ resource.ResourceWithConfigure = (*knowledgeResource)(nil)
)

type knowledgeResource struct {
	name   string
	client *clients.Client
}

func NewKnowledgeResource() resource.Resource {
	return &knowledgeResource{}
}

func (r *knowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_knowledge", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.KnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for knowledge read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log read operation
	logger.Debug("Reading knowledge resource", "knowledge_id", id, "knowledge_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading knowledge resource", sentry.LevelDebug, map[string]interface{}{"knowledge_id": id, "knowledge_name": name})

	if err := r.client.ReadKnowledge(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read knowledge", "knowledge_id", id, "knowledge_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read knowledge", "knowledge_id", id, "knowledge_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *knowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.KnowledgeSchema()
}

func (r *knowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_knowledge", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.KnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for knowledge delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log delete operation
	logger.Info("Deleting knowledge resource", "knowledge_id", id, "knowledge_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting knowledge resource", sentry.LevelInfo, map[string]interface{}{"knowledge_id": id, "knowledge_name": name})

	if err := r.client.DeleteKnowledge(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete knowledge", "knowledge_id", id, "knowledge_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted knowledge", "knowledge_id", id, "knowledge_name", name)
}

func (r *knowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_knowledge", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.KnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for knowledge create", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()

	// Log create operation
	logger.Info("Creating knowledge resource", "knowledge_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating knowledge resource", sentry.LevelInfo, map[string]interface{}{"knowledge_name": name})

	state, err := r.client.CreateKnowledge(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create knowledge", "knowledge_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created knowledge", "knowledge_name", name, "knowledge_id", state.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *knowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_knowledge", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.KnowledgeModel
	var state entities.KnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan or state for knowledge update", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log update operation
	logger.Info("Updating knowledge resource", "knowledge_id", id, "knowledge_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating knowledge resource", sentry.LevelInfo, map[string]interface{}{"knowledge_id": id, "knowledge_name": name})

	updatedState := state

	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		updatedState.Id = plan.Id
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatedState.Name = plan.Name
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		updatedState.Type = plan.Type
	}

	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		updatedState.Groups = plan.Groups
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		updatedState.Labels = plan.Labels
	}

	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		updatedState.Content = plan.Content
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatedState.Description = plan.Description
	}

	if !plan.SupportedAgents.IsNull() && !plan.SupportedAgents.IsUnknown() {
		updatedState.SupportedAgents = plan.SupportedAgents
	}

	if err := r.client.UpdateKnowledge(ctx, &updatedState); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update knowledge", "knowledge_id", id, "knowledge_name", name, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated knowledge", "knowledge_id", id, "knowledge_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *knowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

func (r *knowledgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			logger.Error("Failed to configure knowledge resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "knowledge"
		r.client = client
		logger.Debug("Successfully configured knowledge resource")
	}
}

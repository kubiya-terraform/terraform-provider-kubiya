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
	_ resource.Resource                = (*agentResource)(nil)
	_ resource.ResourceWithConfigure   = (*agentResource)(nil)
	_ resource.ResourceWithImportState = (*agentResource)(nil)
)

type agentResource struct {
	name   string
	client *clients.Client
}

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

func (r *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_agent", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.AgentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for agent read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()

	// Log read operation
	logger.Debug("Reading agent resource", "agent_id", id)
	kubiyasentry.AddBreadcrumb("resource", "Reading agent resource", sentry.LevelDebug, map[string]interface{}{"agent_id": id})

	updatedState, err := r.client.ReadAgent(ctx, id)
	if err != nil || updatedState == nil {
		if err == nil {
			err = fmt.Errorf("agent %s not found", id)
		}
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read agent", "agent_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read agent", "agent_id", id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *agentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.AgentSchema()
}

func (r *agentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_agent", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.AgentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for agent deletion", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log deletion operation
	logger.Info("Deleting agent resource", "agent_id", id, "agent_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting agent resource", sentry.LevelInfo, map[string]interface{}{
		"agent_id":   id,
		"agent_name": name,
	})

	if err := r.client.DeleteAgent(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		CaptureResourceError(ctx, "kubiya_agent", id, "delete", err)
		logger.Error("Failed to delete agent", "agent_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted agent", "agent_id", id, "agent_name", name)
}

func (r *agentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_agent", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.AgentModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for agent creation", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()

	// Log creation operation
	logger.Info("Creating agent resource",
		"agent_name", name,
		"model", plan.Model.ValueString(),
		"runner", plan.Runner.ValueString(),
	)
	kubiyasentry.AddBreadcrumb("resource", "Creating agent resource", sentry.LevelInfo, map[string]interface{}{
		"agent_name": name,
	})

	// Validate agent configuration including MCP server if present
	if err := entities.ValidateAgent(&plan); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Agent Configuration",
			err.Error(),
		)
		return
	}

	state, err := r.client.CreateAgent(ctx, &plan)
	if err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		CaptureResourceError(ctx, "kubiya_agent", "", "create", err)
		logger.Error("Failed to create agent", "agent_name", name, "error", err)

		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	// Success
	id := state.Id.ValueString()
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created agent", "agent_id", id, "agent_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *agentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_agent", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.AgentModel
	var state entities.AgentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan/state for agent update", "error", "diagnostics error")
		return
	}

	updatedState := state

	if plan.Tasks != nil {
		updatedState.Tasks = plan.Tasks
	}

	if plan.Starters != nil {
		updatedState.Starters = plan.Starters
	}

	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		updatedState.Id = plan.Id
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatedState.Name = plan.Name
	}
	if !plan.Image.IsNull() && !plan.Image.IsUnknown() {
		updatedState.Image = plan.Image
	}
	if !plan.Model.IsNull() && !plan.Model.IsUnknown() {
		updatedState.Model = plan.Model
	}
	if !plan.Owner.IsNull() && !plan.Owner.IsUnknown() {
		updatedState.Owner = plan.Owner
	}
	if !plan.Runner.IsNull() && !plan.Runner.IsUnknown() {
		updatedState.Runner = plan.Runner
	}
	if !plan.CreatedAt.IsNull() && !plan.CreatedAt.IsUnknown() {
		updatedState.CreatedAt = plan.CreatedAt
	}
	if !plan.IsDebugMode.IsNull() && !plan.IsDebugMode.IsUnknown() {
		updatedState.IsDebugMode = plan.IsDebugMode
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatedState.Description = plan.Description
	}
	if !plan.Instructions.IsNull() && !plan.Instructions.IsUnknown() {
		updatedState.Instructions = plan.Instructions
	}

	if !plan.Links.IsNull() && !plan.Links.IsUnknown() {
		updatedState.Links = plan.Links
	}
	if !plan.Tools.IsNull() && !plan.Tools.IsUnknown() {
		updatedState.Tools = plan.Tools
	}
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		updatedState.Users = plan.Users
	}
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		updatedState.Groups = plan.Groups
	}
	if !plan.Secrets.IsNull() && !plan.Secrets.IsUnknown() {
		updatedState.Secrets = plan.Secrets
	}
	if !plan.Sources.IsNull() && !plan.Sources.IsUnknown() {
		updatedState.Sources = plan.Sources
	}
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		updatedState.Variables = plan.Variables
	}
	if !plan.Integrations.IsNull() && !plan.Integrations.IsUnknown() {
		updatedState.Integrations = plan.Integrations
	}

	// Handle MCP server updates
	if plan.MCPServer != nil {
		updatedState.MCPServer = plan.MCPServer
	} else if plan.MCPServer == nil && state.MCPServer != nil {
		// If plan has nil MCP server but state has one, user wants to remove it
		updatedState.MCPServer = nil
	}

	// Validate the updated agent configuration including MCP server if present
	if err := entities.ValidateAgent(&updatedState); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Agent Configuration",
			err.Error(),
		)
		return
	}

	id := updatedState.Id.ValueString()
	name := updatedState.Name.ValueString()

	// Log update operation
	logger.Info("Updating agent resource", "agent_id", id, "agent_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating agent resource", sentry.LevelInfo, map[string]interface{}{
		"agent_id":   id,
		"agent_name": name,
	})

	if err := r.client.UpdateAgent(ctx, &updatedState); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		CaptureResourceError(ctx, "kubiya_agent", id, "update", err)
		logger.Error("Failed to update agent", "agent_id", id, "error", err)
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated agent", "agent_id", id, "agent_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *agentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *agentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
	}

	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			logger.Error("Failed to configure agent resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "agent"
		r.client = client
		logger.Debug("Configured agent resource")
	}
}

func (r *agentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

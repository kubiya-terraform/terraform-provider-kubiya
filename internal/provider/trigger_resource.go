package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/getsentry/sentry-go"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"

	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &triggerResource{}
	_ resource.ResourceWithConfigure   = &triggerResource{}
	_ resource.ResourceWithImportState = &triggerResource{}
)

// NewTriggerResource is a helper function to simplify the provider implementation.
func NewTriggerResource() resource.Resource {
	return &triggerResource{}
}

// triggerResource is the resource implementation.
type triggerResource struct {
	client *clients.Client
}

// Metadata returns the resource type name.
func (r *triggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger"
}

// Schema defines the schema for the resource.
func (r *triggerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.TriggerSchema()
}

// Configure adds the provider configured client to the resource.
func (r *triggerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		logger.Error("Failed to configure trigger resource", "error", "invalid provider data type")
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
	logger.Debug("Successfully configured trigger resource")
}

// Create creates the resource and sets the initial Terraform state.
func (r *triggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_trigger", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Retrieve values from plan
	var plan entities.TriggerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for trigger create", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()

	// Log create operation
	logger.Info("Creating trigger resource", "trigger_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating trigger resource", sentry.LevelInfo, map[string]interface{}{"trigger_name": name})

	// Create the trigger
	tflog.Debug(ctx, "Creating trigger", map[string]interface{}{
		"name": name,
	})

	createdTrigger, err := r.client.CreateTrigger(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create trigger", "trigger_name", name, "error", err)
		resp.Diagnostics.AddError(
			"Error creating trigger",
			"Could not create trigger, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Id = createdTrigger.Id
	plan.Url = createdTrigger.Url
	plan.Status = createdTrigger.Status
	plan.WorkflowId = createdTrigger.WorkflowId

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to set state for trigger create", "trigger_name", name, "error", "diagnostics error")
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created trigger", "trigger_name", name, "trigger_id", plan.Id.ValueString())
	tflog.Debug(ctx, "Created trigger", map[string]interface{}{
		"id":   plan.Id.ValueString(),
		"name": plan.Name.ValueString(),
	})
}

// Read refreshes the Terraform state with the latest data.
func (r *triggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_trigger", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	// Get current state
	var state entities.TriggerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for trigger read", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log read operation
	logger.Debug("Reading trigger resource", "trigger_id", id, "trigger_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading trigger resource", sentry.LevelDebug, map[string]interface{}{"trigger_id": id, "trigger_name": name})

	tflog.Debug(ctx, "Reading trigger", map[string]interface{}{
		"id":   id,
		"name": name,
	})

	// Get refreshed trigger value from Kubiya
	err := r.client.ReadTrigger(ctx, &state)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read trigger", "trigger_id", id, "trigger_name", name, "error", err)
		resp.Diagnostics.AddError(
			"Error Reading Kubiya Trigger",
			"Could not read Kubiya trigger ID "+id+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to set state for trigger read", "trigger_id", id, "trigger_name", name, "error", "diagnostics error")
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read trigger", "trigger_id", id, "trigger_name", name)
	tflog.Debug(ctx, "Read trigger", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *triggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_trigger", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	// Retrieve values from plan
	var plan entities.TriggerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan for trigger update", "error", "diagnostics error")
		return
	}

	// Get current state to preserve computed values
	var state entities.TriggerModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for trigger update", "error", "diagnostics error")
		return
	}

	// Preserve the ID and workflow_id from state
	plan.Id = state.Id
	plan.WorkflowId = state.WorkflowId

	id := plan.Id.ValueString()
	name := plan.Name.ValueString()

	// Log update operation
	logger.Info("Updating trigger resource", "trigger_id", id, "trigger_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating trigger resource", sentry.LevelInfo, map[string]interface{}{"trigger_id": id, "trigger_name": name})

	tflog.Debug(ctx, "Updating trigger", map[string]interface{}{
		"id":   id,
		"name": name,
	})

	// Update existing trigger
	err := r.client.UpdateTrigger(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update trigger", "trigger_id", id, "trigger_name", name, "error", err)
		resp.Diagnostics.AddError(
			"Error Updating Kubiya Trigger",
			"Could not update trigger, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated trigger to get latest computed values
	err = r.client.ReadTrigger(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to read updated trigger", "trigger_id", id, "trigger_name", name, "error", err)
		resp.Diagnostics.AddError(
			"Error Reading Updated Kubiya Trigger",
			"Could not read updated trigger, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to set state for trigger update", "trigger_id", id, "trigger_name", name, "error", "diagnostics error")
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated trigger", "trigger_id", id, "trigger_name", name)
	tflog.Debug(ctx, "Updated trigger", map[string]interface{}{
		"id":   plan.Id.ValueString(),
		"name": plan.Name.ValueString(),
	})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *triggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_trigger", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	// Retrieve values from state
	var state entities.TriggerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for trigger delete", "error", "diagnostics error")
		return
	}

	id := state.Id.ValueString()
	name := state.Name.ValueString()

	// Log delete operation
	logger.Info("Deleting trigger resource", "trigger_id", id, "trigger_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting trigger resource", sentry.LevelInfo, map[string]interface{}{"trigger_id": id, "trigger_name": name})

	tflog.Debug(ctx, "Deleting trigger", map[string]interface{}{
		"id":   id,
		"name": name,
	})

	// Delete existing trigger
	err := r.client.DeleteTrigger(ctx, &state)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete trigger", "trigger_id", id, "trigger_name", name, "error", err)
		resp.Diagnostics.AddError(
			"Error Deleting Kubiya Trigger",
			"Could not delete trigger, unexpected error: "+err.Error(),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted trigger", "trigger_id", id, "trigger_name", name)
	tflog.Debug(ctx, "Deleted trigger", map[string]interface{}{
		"id":   id,
		"name": name,
	})
}

// ImportState imports an existing resource into Terraform.
func (r *triggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

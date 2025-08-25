package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
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
func (r *triggerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *triggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan entities.TriggerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the trigger
	tflog.Debug(ctx, "Creating trigger", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createdTrigger, err := r.client.CreateTrigger(ctx, &plan)
	if err != nil {
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
		return
	}

	tflog.Debug(ctx, "Created trigger", map[string]interface{}{
		"id":   plan.Id.ValueString(),
		"name": plan.Name.ValueString(),
	})
}

// Read refreshes the Terraform state with the latest data.
func (r *triggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state entities.TriggerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading trigger", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})

	// Get refreshed trigger value from Kubiya
	err := r.client.ReadTrigger(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kubiya Trigger",
			"Could not read Kubiya trigger ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Read trigger", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *triggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan entities.TriggerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to preserve computed values
	var state entities.TriggerModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the ID and workflow_id from state
	plan.Id = state.Id
	plan.WorkflowId = state.WorkflowId

	tflog.Debug(ctx, "Updating trigger", map[string]interface{}{
		"id":   plan.Id.ValueString(),
		"name": plan.Name.ValueString(),
	})

	// Update existing trigger
	err := r.client.UpdateTrigger(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kubiya Trigger",
			"Could not update trigger, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated trigger to get latest computed values
	err = r.client.ReadTrigger(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Kubiya Trigger",
			"Could not read updated trigger, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated trigger", map[string]interface{}{
		"id":   plan.Id.ValueString(),
		"name": plan.Name.ValueString(),
	})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *triggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state entities.TriggerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting trigger", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})

	// Delete existing trigger
	err := r.client.DeleteTrigger(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kubiya Trigger",
			"Could not delete trigger, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Deleted trigger", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// ImportState imports an existing resource into Terraform.
func (r *triggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
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

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.SourceModel
	var state entities.SourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since we don't have a direct UpdateSource method, we'll use a read-modify approach
	// First, read the current state to ensure we have the latest data
	currentState, err := r.client.ReadSource(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	// Create a new instance with the updated values
	// Preserve the ID and other computed fields
	updatePlan := &entities.SourceModel{
		Id:            currentState.Id,
		Name:          currentState.Name,
		Url:           plan.Url,
		Runner:        plan.Runner,
		DynamicConfig: plan.DynamicConfig,
	}

	// Delete and recreate the resource since we don't have a direct update method
	if err := r.client.DeleteSource(ctx, currentState); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	newState, err := r.client.CreateSource(ctx, updatePlan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	updatedState, err := r.client.ReadSource(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *sourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SourceSchema()
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSource(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.SourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.CreateSource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "source"
		r.client = client
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

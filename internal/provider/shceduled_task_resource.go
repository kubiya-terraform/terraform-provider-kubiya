package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
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
	var state entities.ScheduledTaskModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	updatedState, err := r.client.ReadScheduledTask(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *scheduledTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ScheduledTaskSchema()
}

func (r *scheduledTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.ScheduledTaskModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteScheduledTask(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *scheduledTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.ScheduledTaskModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.CreateScheduledTask(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scheduledTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduled_task"
}

func (r *scheduledTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "scheduled_task"
		r.client = client
	}
}

func (r *scheduledTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

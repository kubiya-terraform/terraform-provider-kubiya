package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource                = (*inlineSourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*inlineSourceResource)(nil)
	_ resource.ResourceWithImportState = (*inlineSourceResource)(nil)
)

type inlineSourceResource struct {
	name   string
	client *clients.Client
}

func NewInlineSourceResource() resource.Resource {
	return &inlineSourceResource{}
}

func (r *inlineSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.InlineSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	updatedState, err := r.client.ReadInlineSource(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *inlineSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.InlineSourceSchema()
}

func (r *inlineSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.InlineSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Tools == nil || len(plan.Tools) == 0 {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, "tools is required"),
		)
		return
	}

	state, err := r.client.CreateInlineSource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *inlineSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.InlineSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteInlineSource(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *inlineSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.InlineSourceModel
	var state entities.InlineSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if plan.Tools == nil || len(plan.Tools) == 0 {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, "tools is required"),
		)
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	if plan.Tools != nil {
		updatedState.Tools = plan.Tools
	}
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		updatedState.Id = plan.Id
	}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatedState.Name = plan.Name
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		updatedState.Type = plan.Type
	}
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		updatedState.Config = plan.Config
	}
	if !plan.Runner.IsNull() && !plan.Runner.IsUnknown() {
		updatedState.Runner = plan.Runner
	}

	if err := r.client.UpdateInlineSource(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *inlineSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inline_source"
}

func (r *inlineSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "inline_source"
		r.client = client
	}
}

func (r *inlineSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*secreResource)(nil)
	_ resource.ResourceWithConfigure = (*secreResource)(nil)
)

type secreResource struct {
	name   string
	client *clients.Client
}

func NewSecreResource() resource.Resource {
	return &secreResource{}
}

func (r *secreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	if err := r.client.ReadSecret(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"secret not found",
			fmt.Sprintf("secret by name: %s not found. Error: ", state.Name)+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *secreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.SecretModel
	var state entities.SecretModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		updatedState.Name = plan.Name
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		updatedState.Description = plan.Description
	}

	if err := r.client.UpdateSecret(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			"failed to update secret. Error: "+err.Error(),
			"failed to update secret",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *secreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SecretSchema()
}

func (r *secreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.SecretModel
	// diags := req.Plan.Get(ctx, &obj)
	fmt.Printf("🆎🆎🆎🆎🆎🆎 create %v 🆎🆎🆎🆎🆎🆎🆎\n", reflect.ValueOf(req.Config.Raw).Type().String())
	req.Config.Get(ctx, &plan)
	// plan = obj.(entities.SecretModel)
	// resp.Diagnostics.Append(diags...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	state, err := r.client.CreateSecret(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create secret",
			"failed to create secret. Error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *secreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	if err := r.client.DeleteSecret(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"failed to delete secret",
			"failed to delete secret. Error: "+err.Error(),
		)
	}
}

func (r *secreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *secreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "secret"
		r.client = client
	}
}

package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*integrationResource)(nil)
	_ resource.ResourceWithConfigure = (*integrationResource)(nil)
)

type integrationResource struct {
	name   string
	client *clients.Client
}

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.IntegrationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[Resource] [Integration] [Read] - State: %v", state)

	//if err := r.client.ReadAgent(ctx, &state); err != nil {
	//	resp.Diagnostics.AddError(
	//		resourceActionError(readAction, r.name, err.Error()),
	//	)
	//	return
	//}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.IntegrationSchema()
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.IntegrationModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[Resource] [Integration] [Delete] - State: %v", state)

	//if err := r.client.DeleteAgent(ctx, &state); err != nil {
	//	resp.Diagnostics.AddError(
	//		resourceActionError(deleteAction, r.name, err.Error()),
	//	)
	//}
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.IntegrationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.CreateIntegration(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	//
	//if len(state.Tasks.ValueString()) == len(plan.Tasks.ValueString()) {
	//	state.Tasks = types.StringValue(plan.Tasks.ValueString())
	//}
	//if len(state.Starters.ValueString()) == len(plan.Starters.ValueString()) {
	//	state.Starters = types.StringValue(plan.Starters.ValueString())
	//}
	//
	//resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.IntegrationModel
	var state entities.IntegrationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	log.Printf("[Resource] [Integration] [Update] - Plan: %v", plan)
	log.Printf("[Resource] [Integration] [Update] - State: %v", state)

	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//updatedState := state
	//
	//if !plan.Email.IsNull() {
	//	updatedState.Email = plan.Email
	//}
	//
	//if !plan.Name.IsNull() {
	//	updatedState.Name = plan.Name
	//}
	//
	//if !plan.Image.IsNull() {
	//	updatedState.Image = plan.Image
	//}
	//
	//if !plan.Model.IsNull() {
	//	updatedState.Model = plan.Model
	//}
	//
	//if !plan.Links.IsNull() {
	//	updatedState.Links = plan.Links
	//}
	//
	//if !plan.Users.IsNull() {
	//	updatedState.Users = plan.Users
	//}
	//
	//if !plan.Groups.IsNull() {
	//	updatedState.Groups = plan.Groups
	//}
	//
	//if !plan.Runners.IsNull() {
	//	updatedState.Runners = plan.Runners
	//}
	//
	//if !plan.Secrets.IsNull() {
	//	updatedState.Secrets = plan.Secrets
	//}
	//
	//if !plan.Starters.IsNull() {
	//	updatedState.Starters = plan.Starters
	//}
	//
	//if !plan.Variables.IsNull() {
	//	updatedState.Variables = plan.Variables
	//}
	//
	//if !plan.Description.IsNull() {
	//	updatedState.Description = plan.Description
	//}
	//
	//if !plan.Instructions.IsNull() {
	//	updatedState.Instructions = plan.Instructions
	//}
	//
	//if !plan.Integrations.IsNull() {
	//	updatedState.Integrations = plan.Integrations
	//}
	//
	//if err := r.client.UpdateAgent(ctx, &updatedState); err != nil {
	//	resp.Diagnostics.AddError(
	//		resourceActionError(updateAction, r.name, err.Error()),
	//	)
	//	return
	//}
	//
	//resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "integration"
		r.client = client
	}
}

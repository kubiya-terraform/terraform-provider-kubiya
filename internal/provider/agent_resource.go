package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*agentResource)(nil)
	_ resource.ResourceWithConfigure = (*agentResource)(nil)
	//_ resource.ResourceWithImportState = (*agentResource)(nil)
)

type agentResource struct {
	name   string
	client *clients.Client
}

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

func (r *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan entities.AgentModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.ReadAgent(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *agentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.AgentSchema()
}

func (r *agentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.AgentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAgent(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *agentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.AgentModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.CreateAgent(ctx, &plan)
	log.Printf("[104]: state: %v, error: %v\n", state, err)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *agentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.AgentModel
	var state entities.AgentModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Id.IsNull() &&
		!state.Id.IsUnknown() {
		plan.Id = state.Id
	}

	if !state.Owner.IsNull() &&
		!state.Owner.IsUnknown() {
		plan.Owner = state.Owner
	}

	if !state.CreatedAt.IsNull() &&
		!state.CreatedAt.IsUnknown() {
		plan.CreatedAt = state.Id
	}

	agent, err := r.client.UpdateAgent(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	// Required
	state.Name = agent.Name
	state.Image = agent.Image
	state.Model = agent.Model
	state.Runner = agent.Runner
	state.Description = agent.Description
	state.Instructions = agent.Instructions

	// Optional
	state.Tasks = agent.Tasks
	state.Links = agent.Links
	state.Users = agent.Users
	state.Groups = agent.Groups
	state.Secrets = agent.Secrets
	state.Starters = agent.Starters
	state.Variables = agent.Variables
	state.Integrations = agent.Integrations

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *agentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *agentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "agent"
		r.client = client
	}
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*agentResource)(nil)
	_ resource.ResourceWithConfigure = (*agentResource)(nil)
)

type agentResource struct {
	name   string
	client *clients.Client
}

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

func (r *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.AgentModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ReadAgent(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state
	updatedState.Tasks = plan.Tasks
	updatedState.Starters = plan.Starters

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
	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		updatedState.Variables = plan.Variables
	}
	if !plan.Integrations.IsNull() && !plan.Integrations.IsUnknown() {
		updatedState.Integrations = plan.Integrations
	}

	if err := r.client.UpdateAgent(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
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

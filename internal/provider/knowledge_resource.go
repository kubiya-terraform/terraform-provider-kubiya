package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*knowledgeResource)(nil)
	_ resource.ResourceWithConfigure = (*knowledgeResource)(nil)
)

type knowledgeResource struct {
	name   string
	client *clients.Client
}

func NewKnowledgeResource() resource.Resource {
	return &knowledgeResource{}
}

func (r *knowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.KnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ReadKnowledge(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *knowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.KnowledgeSchema()
}

func (r *knowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.KnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteKnowledge(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *knowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.KnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.client.CreateKnowledge(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *knowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.KnowledgeModel
	var state entities.KnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		updatedState.Id = plan.Id
	}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updatedState.Name = plan.Name
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		updatedState.Type = plan.Type
	}

	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		updatedState.Groups = plan.Groups
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		updatedState.Labels = plan.Labels
	}

	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		updatedState.Content = plan.Content
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatedState.Description = plan.Description
	}

	if !plan.SupportedAgents.IsNull() && !plan.SupportedAgents.IsUnknown() {
		updatedState.SupportedAgents = plan.SupportedAgents
	}

	if err := r.client.UpdateKnowledge(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *knowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

func (r *knowledgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "knowledge"
		r.client = client
	}
}

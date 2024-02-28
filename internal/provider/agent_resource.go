package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/resource_agent"
)

var (
	_ resource.Resource              = (*agentResource)(nil)
	_ resource.ResourceWithConfigure = (*agentResource)(nil)
)

type agentResource struct {
	client *clients.Client
}

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

func (a *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan resource_agent.AgentModel

	// Read Terraform prior plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	apiReq := plan.Uuid.ValueString()
	apiResp, err := a.client.GetAgentById(apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get agent by id",
			"failed to get agent by id"+err.Error())
	}

	result := resource_agent.AgentModel{
		Id:                   types.StringValue(apiResp.Uuid),
		Name:                 types.StringValue(apiResp.Name),
		Uuid:                 types.StringValue(apiResp.Uuid),
		Image:                types.StringValue(apiResp.Image),
		LlmModel:             types.StringValue(apiResp.LlmModel),
		Description:          types.StringValue(apiResp.Description),
		AiInstructions:       types.StringValue(apiResp.AiInstructions),
		Links:                toListType(&resp.Diagnostics, apiResp.Links...),
		Owners:               toListType(&resp.Diagnostics, apiResp.Owners...),
		Runners:              toListType(&resp.Diagnostics, apiResp.Runners...),
		Secrets:              toListType(&resp.Diagnostics, apiResp.Secrets...),
		Starters:             toListType(&resp.Diagnostics, apiResp.Starters...),
		AllowedUsers:         toListType(&resp.Diagnostics, apiResp.AllowedUsers...),
		Integrations:         toListType(&resp.Diagnostics, apiResp.Integrations...),
		AllowedGroups:        toListType(&resp.Diagnostics, apiResp.AllowedGroups...),
		EnvironmentVariables: convertStringMapToMapType(&resp.Diagnostics, apiResp.EnvironmentVariables),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)

	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *agentResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_agent.AgentResourceSchema(ctx)
}

func (a *agentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan resource_agent.AgentModel

	// Read Terraform prior plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	deleteRequest := plan.Uuid.ValueString()
	err := a.client.DeleteAgent(deleteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete agent by id",
			"failed to delete agent by id"+err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (a *agentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_agent.AgentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	apiReq := clients.Agent{
		Name:           plan.Name.ValueString(),
		Image:          plan.Image.ValueString(),
		LlmModel:       plan.LlmModel.ValueString(),
		Description:    plan.Description.ValueString(),
		AiInstructions: plan.AiInstructions.ValueString(),

		Links:                toStringSlice(plan.Links),
		Owners:               toStringSlice(plan.Owners),
		Runners:              toStringSlice(plan.Runners),
		Secrets:              toStringSlice(plan.Secrets),
		Starters:             toStringSlice(plan.Starters),
		AllowedUsers:         toStringSlice(plan.AllowedUsers),
		Integrations:         toStringSlice(plan.Integrations),
		AllowedGroups:        toStringSlice(plan.AllowedGroups),
		EnvironmentVariables: convertTypesMapToStringMap(plan.EnvironmentVariables),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := a.client.CreateAgent(&apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Agent",
			"failed to create agent. Error: "+err.Error(),
		)
		return
	}

	result := resource_agent.AgentModel{
		Id:                   types.StringValue(apiResp.Uuid),
		Name:                 types.StringValue(apiResp.Name),
		Uuid:                 types.StringValue(apiResp.Uuid),
		Image:                types.StringValue(apiResp.Image),
		LlmModel:             types.StringValue(apiResp.LlmModel),
		Description:          types.StringValue(apiResp.Description),
		AiInstructions:       types.StringValue(apiResp.AiInstructions),
		Links:                toListType(&resp.Diagnostics, apiResp.Links...),
		Owners:               toListType(&resp.Diagnostics, apiResp.Owners...),
		Runners:              toListType(&resp.Diagnostics, apiResp.Runners...),
		Secrets:              toListType(&resp.Diagnostics, apiResp.Secrets...),
		Starters:             toListType(&resp.Diagnostics, apiResp.Starters...),
		AllowedUsers:         toListType(&resp.Diagnostics, apiResp.AllowedUsers...),
		Integrations:         toListType(&resp.Diagnostics, apiResp.Integrations...),
		AllowedGroups:        toListType(&resp.Diagnostics, apiResp.AllowedGroups...),
		EnvironmentVariables: convertStringMapToMapType(&resp.Diagnostics, apiResp.EnvironmentVariables),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (a *agentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_agent.AgentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	apiReq := clients.Agent{
		Name:           plan.Name.ValueString(),
		Image:          plan.Image.ValueString(),
		LlmModel:       plan.LlmModel.ValueString(),
		Description:    plan.Description.ValueString(),
		AiInstructions: plan.AiInstructions.ValueString(),

		Links:                toStringSlice(plan.Links),
		Owners:               toStringSlice(plan.Owners),
		Runners:              toStringSlice(plan.Runners),
		Secrets:              toStringSlice(plan.Secrets),
		Starters:             toStringSlice(plan.Starters),
		AllowedUsers:         toStringSlice(plan.AllowedUsers),
		Integrations:         toStringSlice(plan.Integrations),
		AllowedGroups:        toStringSlice(plan.AllowedGroups),
		EnvironmentVariables: convertTypesMapToStringMap(plan.EnvironmentVariables),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := a.client.CreateAgent(&apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Agent",
			"failed to create agent. Error: "+err.Error(),
		)
		return
	}

	result := resource_agent.AgentModel{
		Id:                   types.StringValue(apiResp.Uuid),
		Name:                 types.StringValue(apiResp.Name),
		Uuid:                 types.StringValue(apiResp.Uuid),
		Image:                types.StringValue(apiResp.Image),
		LlmModel:             types.StringValue(apiResp.LlmModel),
		Description:          types.StringValue(apiResp.Description),
		AiInstructions:       types.StringValue(apiResp.AiInstructions),
		Links:                toListType(&resp.Diagnostics, apiResp.Links...),
		Owners:               toListType(&resp.Diagnostics, apiResp.Owners...),
		Runners:              toListType(&resp.Diagnostics, apiResp.Runners...),
		Secrets:              toListType(&resp.Diagnostics, apiResp.Secrets...),
		Starters:             toListType(&resp.Diagnostics, apiResp.Starters...),
		AllowedUsers:         toListType(&resp.Diagnostics, apiResp.AllowedUsers...),
		Integrations:         toListType(&resp.Diagnostics, apiResp.Integrations...),
		AllowedGroups:        toListType(&resp.Diagnostics, apiResp.AllowedGroups...),
		EnvironmentVariables: convertStringMapToMapType(&resp.Diagnostics, apiResp.EnvironmentVariables),
	}

	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (a *agentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (a *agentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *clients.AgentsClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.client = client
}

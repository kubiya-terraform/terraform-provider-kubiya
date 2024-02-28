package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/resource_agent"
	"terraform-provider-kubiya/internal/resource_runner"
)

var (
	_ resource.Resource              = (*runnerResource)(nil)
	_ resource.ResourceWithConfigure = (*runnerResource)(nil)
)

type runnerResource struct {
	client *clients.Client
}

func NewRunnerResource() resource.Resource {
	return &runnerResource{}
}

func (r *runnerResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	//var data resource_runner.RunnerModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//// Read API call logic
	//
	//// Save updated data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *runnerResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *runnerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_agent.AgentResourceSchema(ctx)
}

func (r *runnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_runner.RunnerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	name := data.Name.ValueString()
	runner, err := r.client.CreateRunner(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create runner",
			"failed to create runner"+err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &runner)...)
}

func (r *runnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_runner.RunnerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	name := data.Name.ValueString()
	err := r.client.DeleteRunner(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete runner",
			"failed to delete runner"+err.Error(),
		)
	}
}

func (r *runnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner"
}

func (r *runnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
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

func (r *runnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_runner.RunnerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *runnerResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *runnerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_runner.RunnerResourceSchema(ctx)
}

func (r *runnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_runner.RunnerModel

	var validator = func(name string) bool {
		// Updated the regex pattern to enforce a minimum length of 5 characters
		pattern := `^[a-z][a-z0-9-]{4,}[a-z]$`
		regex := regexp.MustCompile(pattern)

		// Test the input string against the pattern
		return regex.MatchString(name)
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	name := data.Name.ValueString()
	if valid := validator(name); !valid {
		resp.Diagnostics.AddError(
			"Name is not valid name",
			"Runner name can only contain lowercase alpha-numeric characters. Spaces and underscores are not allowed.",
		)
		return
	}

	path := data.RunnerDeploymentFolder.ValueString()
	if data.RunnerDeploymentFolder.IsNull() ||
		data.RunnerDeploymentFolder.IsUnknown() {
		resp.Diagnostics.AddError(
			"filepath is missing or empty",
			"filepath is missing or empty. ex. /Users/mevrat.avraham/, /Users/mevrat.avraham/runners/",
		)
		return
	}

	runner, err := r.client.CreateRunner(name, path)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create runner",
			"failed to create runner. Error: "+err.Error(),
		)
		return
	}

	data.Url = types.StringValue(runner.Url)
	data.RunnerDeploymentFile = types.StringValue(runner.Path)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
			"failed to delete runner. Error: "+err.Error(),
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

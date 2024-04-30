package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
)

var (
	_ resource.Resource              = (*webhookResource)(nil)
	_ resource.ResourceWithConfigure = (*webhookResource)(nil)
)

type webhookModel struct {
	Id            types.String        `tfsdk:"id"`
	Org           types.String        `tfsdk:"org"`
	Name          types.String        `tfsdk:"name"`
	Source        types.String        `tfsdk:"source"`
	Prompt        types.String        `tfsdk:"prompt"`
	Filter        types.String        `tfsdk:"filter"`
	AgentId       types.String        `tfsdk:"agent_id"`
	CreatedBy     types.String        `tfsdk:"created_by"`
	Webhook       types.String        `tfsdk:"webhook_url"`
	Communication *communicationModel `tfsdk:"communication"`
}

type communicationModel struct {
	Method      string `tfsdk:"method"`
	Destination string `tfsdk:"destination"`
}

type webhookResource struct {
	client *clients.Client
}

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var d entities.WebhookModel
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	result, err := r.client.GetWebhook(&d)
	if err != nil || result == nil {
		if err != nil {
			diags.AddError("webhook not found", err.Error())
		} else {
			diags.AddError("webhook not found",
				fmt.Sprintf("webhook by name: %s not found", d.Name),
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entities.WebhookModel
	var state entities.WebhookModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	if !plan.Id.IsUnknown() && !plan.Id.IsNull() {
		updatedState.Id = plan.Id
	}
	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		updatedState.Name = plan.Name
	}
	if !plan.Agent.IsUnknown() && !plan.Agent.IsNull() {
		updatedState.Agent = plan.Agent
	}
	if !plan.Prompt.IsUnknown() && !plan.Prompt.IsNull() {
		updatedState.Prompt = plan.Prompt
	}
	if !plan.Source.IsUnknown() && !plan.Source.IsNull() {
		updatedState.Source = plan.Source
	}
	if !plan.Filter.IsUnknown() && !plan.Filter.IsNull() {
		updatedState.Filter = plan.Filter
	}
	if !plan.Destination.IsUnknown() && !plan.Destination.IsNull() {
		updatedState.Destination = plan.Destination
	}

	result, err := r.client.UpdateWebhook(&updatedState)
	if err != nil {
		resp.Diagnostics.AddError("failed to update webhook", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.WebhookSchema()
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entities.WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateWebhook(&plan)
	if err != nil {
		resp.Diagnostics.AddError("failed to update webhook", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data entities.WebhookModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	if err := r.client.DeleteWebhook(&data); err != nil {
		resp.Diagnostics.AddError(
			"failed to delete webhook",
			"failed to delete webhook. Error: "+err.Error(),
		)
	}
}

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

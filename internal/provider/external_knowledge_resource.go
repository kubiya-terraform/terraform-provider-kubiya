package provider

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

var (
	_ resource.Resource              = (*externalKnowledgeResource)(nil)
	_ resource.ResourceWithConfigure = (*externalKnowledgeResource)(nil)
)

type externalKnowledgeResource struct {
	name   string
	client *clients.Client
}

func NewExternalKnowledgeResource() resource.Resource {
	return &externalKnowledgeResource{}
}

func (r *externalKnowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Reading external_knowledge resource", sentry.LevelInfo, nil)

	var state entities.ExternalKnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ReadExternalKnowledge(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *externalKnowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ExternalKnowledgeSchema()
}

func (r *externalKnowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Deleting external_knowledge resource", sentry.LevelInfo, nil)

	var state entities.ExternalKnowledgeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteExternalKnowledge(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *externalKnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Creating external_knowledge resource", sentry.LevelInfo, nil)

	var plan entities.ExternalKnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	state, err := r.client.CreateExternalKnowledge(ctx, &plan)
	if err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_external_knowledge", "", "create", err)

		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *externalKnowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_external_knowledge", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Updating external_knowledge resource", sentry.LevelInfo, nil)

	var plan entities.ExternalKnowledgeModel
	var state entities.ExternalKnowledgeModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := state

	// Update vendor if it has changed
	if !plan.Vendor.IsNull() && !plan.Vendor.IsUnknown() {
		updatedState.Vendor = plan.Vendor
	}

	// Update config if it has changed
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		updatedState.Config = plan.Config
	}

	if err := r.client.UpdateExternalKnowledge(ctx, &updatedState); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(updateAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *externalKnowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_knowledge"
}

func (r *externalKnowledgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "external_knowledge"
		r.client = client
	}
}

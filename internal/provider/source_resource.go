package provider

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

var (
	_ resource.Resource                = (*sourceResource)(nil)
	_ resource.ResourceWithConfigure   = (*sourceResource)(nil)
	_ resource.ResourceWithImportState = (*sourceResource)(nil)
)

type sourceResource struct {
	name   string
	client *clients.Client
}

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

func (r *sourceResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	updatedState, err := r.client.ReadSource(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *sourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SourceSchema()
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state entities.SourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSource(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_source", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Creating source resource", sentry.LevelInfo, nil)

	var plan entities.SourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	state, err := r.client.CreateSource(ctx, &plan)
	if err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_source", "", "create", err)

		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "source"
		r.client = client
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

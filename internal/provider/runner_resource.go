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
	_ resource.Resource              = (*runnerResource)(nil)
	_ resource.ResourceWithConfigure = (*runnerResource)(nil)
)

type runnerResource struct {
	name   string
	client *clients.Client
}

func NewRunnerResource() resource.Resource {
	return &runnerResource{}
}

func (r *runnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Reading runner resource", sentry.LevelInfo, nil)

	var state entities.RunnerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Update span with resource name if available
	name := ""
	if !state.Name.IsNull() {
		name = state.Name.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, name)
		kubiyasentry.SetSpanData(span, "runner.name", name)
	}

	if err := r.client.ReadRunner(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_runner", name, "read", err)

		resp.Diagnostics.AddError(
			resourceActionError(readAction, r.name, err.Error()),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *runnerResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *runnerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.RunnerSchema()
}

func (r *runnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Creating runner resource", sentry.LevelInfo, nil)

	var plan entities.RunnerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Add plan data to span
	if !plan.Name.IsNull() {
		kubiyasentry.SetSpanData(span, "runner.name", plan.Name.ValueString())
	}

	state, err := r.client.CreateRunner(ctx, &plan)
	if err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_runner", "", "create", err)

		resp.Diagnostics.AddError(
			resourceActionError(createAction, r.name, err.Error()),
		)
		return
	}

	// Update span with created resource name
	if state != nil && !state.Name.IsNull() {
		name := state.Name.ValueString()
		kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, name)
		kubiyasentry.SetSpanData(span, "runner.created_name", name)
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *runnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_runner", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Deleting runner resource", sentry.LevelInfo, nil)

	var state entities.RunnerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Update span with resource name
	name := state.Name.ValueString()
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, name)
	kubiyasentry.SetSpanData(span, "runner.name", name)

	if err := r.client.DeleteRunner(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_runner", name, "delete", err)

		resp.Diagnostics.AddError(
			resourceActionError(deleteAction, r.name, err.Error()),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	kubiyasentry.AddBreadcrumb("resource", "Successfully deleted runner "+name, sentry.LevelInfo, nil)
}

func (r *runnerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner"
}

func (r *runnerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "runner"
		r.client = client
	}
}

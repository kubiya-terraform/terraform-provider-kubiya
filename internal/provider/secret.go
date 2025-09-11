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
	_ resource.Resource              = (*secretResource)(nil)
	_ resource.ResourceWithConfigure = (*secretResource)(nil)
)

type secretResource struct {
	name   string
	client *clients.Client
}

func NewSecreResource() resource.Resource {
	return &secretResource{}
}

func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Reading secret resource", sentry.LevelInfo, nil)

	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Read API call logic
	if err := r.client.ReadSecret(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_secret", "", "read", err)

		resp.State.RemoveResource(ctx)
		// resp.Diagnostics.AddError(
		// 	"secret not found",
		// 	fmt.Sprintf("secret by name: %s not found. Error: ", state.Name)+err.Error(),
		// )
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Updating secret resource", sentry.LevelInfo, nil)

	var plan entities.SecretModel
	var state entities.SecretModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	updatedState := state

	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		updatedState.Name = plan.Name
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		updatedState.Description = plan.Description
	}
	if !plan.Value.IsUnknown() && !plan.Value.IsNull() {
		updatedState.Value = plan.Value
	}

	if err := r.client.UpdateSecret(ctx, &updatedState); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_secret", "", "update", err)

		resp.Diagnostics.AddError(
			"failed to update secret. Error: "+err.Error(),
			"failed to update secret",
		)
	}

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SecretSchema()
}

func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Creating secret resource", sentry.LevelInfo, nil)

	var plan entities.SecretModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	state, err := r.client.CreateSecret(ctx, &plan)
	if err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_secret", "", "create", err)

		resp.Diagnostics.AddError(
			"failed to create secret",
			"failed to create secret. Error: "+err.Error(),
		)
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("resource", "Deleting secret resource", sentry.LevelInfo, nil)

	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Delete API call logic
	if err := r.client.DeleteSecret(ctx, &state); err != nil {
		// Record error in span and capture to Sentry
		kubiyasentry.RecordError(ctx, err)
		CaptureResourceError(ctx, "kubiya_secret", "", "delete", err)
		// Don't add error to diagnostics as the original code doesn't
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	resp.State.RemoveResource(ctx)
}

func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "secret"
		r.client = client
	}
}

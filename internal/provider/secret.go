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
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the read operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for secret read", "error", "diagnostics error")
		return
	}

	name := state.Name.ValueString()

	// Log read operation
	logger.Debug("Reading secret resource", "secret_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Reading secret resource", sentry.LevelDebug, map[string]interface{}{"secret_name": name})

	// Read API call logic
	if err := r.client.ReadSecret(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusNotFound)
		logger.Error("Failed to read secret", "secret_name", name, "error", err)
		resp.State.RemoveResource(ctx)
		// resp.Diagnostics.AddError(
		// 	"secret not found",
		// 	fmt.Sprintf("secret by name: %s not found. Error: ", state.Name)+err.Error(),
		// )
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Debug("Successfully read secret", "secret_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the update operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.SecretModel
	var state entities.SecretModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get plan or state for secret update", "error", "diagnostics error")
		return
	}

	name := state.Name.ValueString()

	// Log update operation
	logger.Info("Updating secret resource", "secret_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Updating secret resource", sentry.LevelInfo, map[string]interface{}{"secret_name": name})

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
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to update secret", "secret_name", name, "error", err)
		resp.Diagnostics.AddError(
			"failed to update secret. Error: "+err.Error(),
			"failed to update secret",
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully updated secret", "secret_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.SecretSchema()
}

func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the create operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	var plan entities.SecretModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get config for secret create", "error", "diagnostics error")
		return
	}

	name := plan.Name.ValueString()

	// Log create operation
	logger.Info("Creating secret resource", "secret_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Creating secret resource", sentry.LevelInfo, map[string]interface{}{"secret_name": name})

	state, err := r.client.CreateSecret(ctx, &plan)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to create secret", "secret_name", name, "error", err)
		resp.Diagnostics.AddError(
			"failed to create secret",
			"failed to create secret. Error: "+err.Error(),
		)
		return
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully created secret", "secret_name", name)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Start tracing for the delete operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "kubiya_secret", "", kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	var state entities.SecretModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		logger.Error("Failed to get state for secret delete", "error", "diagnostics error")
		return
	}

	name := state.Name.ValueString()

	// Log delete operation
	logger.Info("Deleting secret resource", "secret_name", name)
	kubiyasentry.AddBreadcrumb("resource", "Deleting secret resource", sentry.LevelInfo, map[string]interface{}{"secret_name": name})

	// Delete API call logic
	if err := r.client.DeleteSecret(ctx, &state); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		logger.Error("Failed to delete secret", "secret_name", name, "error", err)
		// Don't add error to diagnostics as the original code doesn't
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	logger.Info("Successfully deleted secret", "secret_name", name)
	resp.State.RemoveResource(ctx)
}

func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *secretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Get or create logger in context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if req.ProviderData != nil {
		var ok bool
		var client *clients.Client

		if client, ok = req.ProviderData.(*clients.Client); !ok {
			logger.Error("Failed to configure secret resource", "error", "invalid provider data type")
			resp.Diagnostics.AddError(configResourceError(req.ProviderData))
			return
		}

		r.name = "secret"
		r.client = client
		logger.Debug("Successfully configured secret resource")
	}
}

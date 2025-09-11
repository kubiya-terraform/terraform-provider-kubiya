package provider

import (
	"context"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type kubiyaProvider struct {
	version string
}

var _ provider.Provider = (*kubiyaProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kubiyaProvider{version: version}
	}
}

func (p *kubiyaProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentResource,
		NewRunnerResource,
		NewSourceResource,
		NewWebhookResource,
		NewKnowledgeResource,
		NewExternalKnowledgeResource,
		NewIntegrationResource,
		NewScheduledTaskResource,
		NewSecreResource,
		NewInlineSourceResource,
		NewTriggerResource,
	}
}

func (p *kubiyaProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *kubiyaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *kubiyaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubiya"
}

func (p *kubiyaProvider) Configure(ctx context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Initialize Sentry when provider is configured
	if err := kubiyasentry.Initialize(p.version); err != nil {
		// Log error but don't fail provider configuration
		// Sentry is optional for functionality
	}

	// Start a transaction for provider configuration
	ctx, span := kubiyasentry.StartTransaction(ctx, "provider.configure", "Configuring Kubiya provider")
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("provider", "Configuring Kubiya provider", sentry.LevelInfo, nil)

	const (
		apiKeyEnvVar         = "KUBIYA_API_KEY"
		envKeyEnvVar         = "KUBIYA_ENV" // New environment variable for environment selection
		missingAPIKey        = "Kubiya API Key Not Configured"
		missingAPIKeyDetails = "Please set the Kubiya API Key using the environment variable 'KUBIYA_API_KEY'. " +
			"Use the command below:\n> export KUBIYA_API_KEY=YOUR_API_KEY"
	)

	apiKey := os.Getenv(apiKeyEnvVar)
	if apiKey == "" {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		resp.Diagnostics.AddError(missingAPIKey, missingAPIKeyDetails)
		return
	}

	// Fetch the environment or set to default
	env := os.Getenv(envKeyEnvVar)
	if env == "" {
		env = "production"
	}

	// Set Sentry environment tag
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("environment", env)
		scope.SetTag("provider.version", p.version)
	})

	// Create a new Kubiya client using the API key and environment
	client, err := clients.New(apiKey, env)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Failed to Create Kubiya Client", "An error occurred while creating the Kubiya client: "+err.Error())
		return
	}

	// Success
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)

	// Attach the client to be used by resources and data sources
	resp.ResourceData = client
	resp.DataSourceData = client
}

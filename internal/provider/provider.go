package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-kubiya/internal/clients"
)

type kubiyaProvider struct {
	version string
}

var _ provider.Provider = (*kubiyaProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kubiyaProvider{
			version: version,
		}
	}
}

func (p *kubiyaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentResource,
		NewRunnerResource,
		NewWebhookResource,
		NewIntegrationResource,
	}
}

func (p *kubiyaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *kubiyaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *kubiyaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubiya"
}

func (p *kubiyaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	const (
		env     = "KUBIYA_API_KEY"
		summery = "Kubiya API Key Not Configured"
		details = "Please set the Kubiya API Key using the environment variable 'KUBIYA_API_KEY'. Use the command below:\n> export KUBIYA_API_KEY=YOUR_API_KEY"
	)

	apiKey := os.Getenv(env)
	if len(apiKey) <= 0 {
		resp.Diagnostics.AddError(summery, details)
		return
	}

	// Create a new Kubiya client using the configuration values
	client, err := clients.New(apiKey)
	if err != nil {
		resp.Diagnostics.AddError(summery, details)
		return
	}

	// Make the Kubiya client available during DataSource and Resource
	// type Configure methods.
	resp.ResourceData = client
	resp.DataSourceData = client
}

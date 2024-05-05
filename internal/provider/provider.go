package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
	"terraform-provider-kubiya/internal/entities"
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
	}
}

func (p *kubiyaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *kubiyaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = entities.ProviderConfigSchema()
}

func (p *kubiyaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubiya"
}

func (p *kubiyaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg entities.ProviderConfig
	diags := req.Config.Get(ctx, &cfg)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.UserKey.IsNull() ||
		cfg.UserKey.IsUnknown() {
		const (
			attr    = "user_key"
			env     = "KUBIYA_USER_KEY"
			summery = "Unknown Kubiya user_key"
			details = "The provider cannot create the Kubiya API client as there is an unknown configuration value for the Kubiya API user_key. Either target apply the source of the value first, set the value statically in the configuration, or use the KUBIYA_USER_KEY environment variable."
		)
		if v := os.Getenv(env); len(v) >= 1 {
			cfg.UserKey = types.StringValue(v)
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root(attr), summery, details)

			return
		}
	}

	// Create a new Kubiya client using the configuration values
	client, err := clients.New(cfg.UserKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(configureProviderError(err))
		return
	}

	// Make the Kubiya client available during DataSource and Resource
	// type Configure methods.
	resp.ResourceData = client
	resp.DataSourceData = client
}

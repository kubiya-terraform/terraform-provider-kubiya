package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/clients"
)

type kubiyaProvider struct {
	version string
}

type KubiyaProviderModel struct {
	Email        types.String `tfsdk:"email"`
	UserKey      types.String `tfsdk:"user_key"`
	Organization types.String `tfsdk:"organization"`
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
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Optional:  true,
				Sensitive: false,
			},
			"user_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: false,
			},
			"organization": schema.StringAttribute{
				Optional:  true,
				Sensitive: false,
			},
		},
	}
}

func (p *kubiyaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "Kubiya"
}

func (p *kubiyaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	const (
		userKeyAttributeName   = "user_key"
		userKeyEnvironmentName = "KUBIYA_USER_KEY"

		apiClientErrSummery = "Unable to Create Kubiya API AgentsClient"
		apiClientErrDetails = "An unexpected error occurred when creating the Kubiya API client. If the error is not clear, please contact the provider developers.\n\nKubiya AgentsClient Error: %s"

		userKeySummery = "Unknown Kubiya UserKey"
		userKeyDetails = "The provider cannot create the Kubiya API client as there is an unknown configuration value for the Kubiya API user_key. Either target apply the source of the value first, set the value statically in the configuration, or use the KUBIYA_USER_KEY environment variable."
	)

	var config KubiyaProviderModel
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.UserKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(userKeyAttributeName),
			userKeySummery, userKeyDetails,
		)
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	userKey := os.Getenv(userKeyEnvironmentName)

	if !config.UserKey.IsNull() {
		userKey = config.UserKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if userKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(userKeyAttributeName),
			userKeySummery, userKeyDetails,
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Kubiya client using the configuration values
	client, err := clients.
		NewClient(userKey,
			config.Email.ValueString(),
			config.Organization.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			apiClientErrSummery,
			fmt.Sprintf(apiClientErrDetails, err.Error()),
		)
		return
	}

	// Make the Kubiya client available during DataSource and Resource
	// type Configure methods.
	resp.ResourceData = client
	resp.DataSourceData = client
}

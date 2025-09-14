package sentry

import "time"

// DSN will be set at build time via ldflags
// Use: -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=your-dsn-here"
var DSN string

// Version will be set at build time via ldflags
// Use: -ldflags "-X terraform-provider-kubiya/internal/sentry.Version=your-dsn-here"
var Version string

const (
	// DefaultDSN is a fallback for local development (empty means Sentry disabled)
	DefaultDSN = ""

	// Default environment if KUBIYA_ENV is not set
	DefaultEnvironment = "production"

	// Supported environment values
	EnvironmentStaging    = "staging"
	EnvironmentProduction = "production"

	// Flush timeout for Sentry on shutdown
	FlushTimeout = 2 * time.Second

	// Transaction operation names
	OpTerraformApply   = "terraform.apply"
	OpTerraformPlan    = "terraform.plan"
	OpTerraformDestroy = "terraform.destroy"
	OpTerraformRefresh = "terraform.refresh"
	OpTerraformImport  = "terraform.import"

	// Span operation names for resources
	OpResourceCreate = "resource.create"
	OpResourceRead   = "resource.read"
	OpResourceUpdate = "resource.update"
	OpResourceDelete = "resource.delete"

	// Span operation names for API calls
	OpAPICall     = "api.call"
	OpStateManage = "state.manage"
	OpValidation  = "validation"

	// Context keys for trace propagation
	TraceHeader   = "sentry-trace"
	BaggageHeader = "baggage"

	// Tags and context keys
	TagResourceType    = "resource.type"
	TagResourceID      = "resource.id"
	TagOperation       = "operation"
	TagProviderVersion = "provider.version"
	TagHTTPMethod      = "http.method"
	TagHTTPURL         = "http.url"
	TagErrorType       = "error.type"
	TagRetryCount      = "retry.count"

	// Breadcrumb types
	BreadcrumbTypeDebug   = "debug"
	BreadcrumbTypeInfo    = "info"
	BreadcrumbTypeHTTP    = "http"
	BreadcrumbTypeError   = "error"
	BreadcrumbTypeDefault = "default"
)

// Sensitive field patterns to redact
var SensitiveFieldPatterns = []string{
	"password",
	"token",
	"secret",
	"key",
	"authorization",
	"api_key",
	"access_token",
	"refresh_token",
	"private_key",
	"client_secret",
}

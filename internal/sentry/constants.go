package sentry

import "time"

// DSN will be set at build time via ldflags
// Use: -ldflags "-X terraform-provider-kubiya/internal/sentry.DSN=your-dsn-here"
var DSN string

const (
	// DefaultDSN is a fallback for local development (empty means Sentry disabled)
	DefaultDSN = ""

	// Default environment if KUBIYA_ENV is not set
	DefaultEnvironment = "production"

	// Supported environment values
	EnvironmentStaging    = "staging"
	EnvironmentProduction = "production"

	// Sample rates for different environments
	// Production has lower sample rates to reduce overhead
	ProductionErrorSampleRate    = 1.0  // Capture all errors in production
	ProductionTracesSampleRate   = 0.1  // 10% of transactions
	ProductionProfilesSampleRate = 0.01 // 1% profiling

	// Staging has higher sample rates for better visibility
	StagingErrorSampleRate    = 1.0 // Capture all errors in staging
	StagingTracesSampleRate   = 0.5 // 50% of transactions
	StagingProfilesSampleRate = 0.1 // 10% profiling

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
	OpAPIAuth     = "api.auth"
	OpAPIRetry    = "api.retry"
	OpStateManage = "state.manage"
	OpValidation  = "validation"

	// Context keys for trace propagation
	TraceHeader   = "sentry-trace"
	BaggageHeader = "baggage"

	// Tags and context keys
	TagResourceType     = "resource.type"
	TagResourceID       = "resource.id"
	TagOperation        = "operation"
	TagTerraformVersion = "terraform.version"
	TagProviderVersion  = "provider.version"
	TagEnvironment      = "environment"
	TagOrganizationID   = "organization.id"
	TagUserID           = "user.id"
	TagHTTPMethod       = "http.method"
	TagHTTPURL          = "http.url"
	TagHTTPStatusCode   = "http.status_code"
	TagErrorType        = "error.type"
	TagRetryCount       = "retry.count"

	// Breadcrumb types
	BreadcrumbTypeDebug      = "debug"
	BreadcrumbTypeInfo       = "info"
	BreadcrumbTypeNavigation = "navigation"
	BreadcrumbTypeHTTP       = "http"
	BreadcrumbTypeError      = "error"
	BreadcrumbTypeDefault    = "default"

	// Performance monitoring
	EnableTracing    = true
	EnableProfiling  = true
	AttachStacktrace = true
	Debug            = false
)

// Resource types for tracking
var ResourceTypes = []string{
	"kubiya_agent",
	"kubiya_runner",
	"kubiya_webhook",
	"kubiya_trigger",
	"kubiya_scheduled_task",
}

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

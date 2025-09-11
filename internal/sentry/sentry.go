package sentry

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

// Config holds the Sentry configuration
type Config struct {
	DSN              string
	Environment      string
	SampleRate       float64
	TracesSampleRate float64
	EnableTracing    bool
	AttachStacktrace bool
	Debug            bool
	ServerName       string
	Release          string
}

// Initialize initializes the Sentry SDK with enhanced configuration
func Initialize() error {
	config := getConfig()

	// Skip initialization if DSN is empty
	if config.DSN == "" {
		// Sentry is disabled when DSN is not provided
		return nil
	}

	// Initialize logging first
	InitializeLogging()
	logger := GetLogger()

	// Use HTTPSyncTransport for synchronous sending
	transport := sentry.NewHTTPSyncTransport()
	logger.Info("Using HTTPSyncTransport")

	clientOptions := sentry.ClientOptions{
		Dsn:              config.DSN,
		SampleRate:       config.SampleRate,
		TracesSampleRate: config.TracesSampleRate,
		EnableTracing:    config.EnableTracing,
		AttachStacktrace: config.AttachStacktrace,
		Debug:            config.Debug,
		ServerName:       config.ServerName,
		Release:          config.Release,
		Transport:        transport, // Always set the transport
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Sanitize sensitive data before sending
			return sanitizeEvent(event)
		},
		BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Sanitize transaction data
			return sanitizeEvent(event)
		},
		BeforeBreadcrumb: func(breadcrumb *sentry.Breadcrumb, hint *sentry.BreadcrumbHint) *sentry.Breadcrumb {
			// Sanitize breadcrumb data
			return sanitizeBreadcrumb(breadcrumb)
		},
		Integrations: func(integrations []sentry.Integration) []sentry.Integration {
			// Remove default integrations that might leak sensitive data
			filtered := make([]sentry.Integration, 0)
			for _, integration := range integrations {
				name := integration.Name()
				// Keep only safe integrations
				if name != "Modules" && name != "Environment" {
					filtered = append(filtered, integration)
				}
			}
			return filtered
		},
	}

	err := sentry.Init(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to initialize Sentry: %w", err)
	}

	// Set initial scope data
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(TagProviderVersion, Version)
		scope.SetContext("runtime", map[string]interface{}{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
		})
	})

	return nil
}

// getConfig returns the Sentry configuration based on environment
func getConfig() *Config {

	// Use build-time DSN if available, otherwise use default (empty = disabled)
	dsn := DSN
	if dsn == "" {
		dsn = DefaultDSN
	}

	config := &Config{
		DSN:              dsn,
		EnableTracing:    true,
		AttachStacktrace: true,
		Debug:            false,
		SampleRate:       1.0,
		TracesSampleRate: 0.1,
		Release:          fmt.Sprintf("terraform-provider-kubiya@%s", Version),
	}

	return config
}

// Flush flushes any buffered events to Sentry
func Flush() {
	sentry.Flush(FlushTimeout)
}

// CaptureError captures an error with additional context
func CaptureError(err error, ctx context.Context, tags map[string]string) {
	if err == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		// Add tags
		for key, value := range tags {
			scope.SetTag(key, value)
		}

		// Add context if available
		if ctx != nil {
			if traceID := GetTraceID(ctx); traceID != "" {
				scope.SetTag("trace_id", traceID)
			}
		}

		sentry.CaptureException(err)
	})
}

// CaptureMessage captures a message with a specific level
func CaptureMessage(message string, level sentry.Level, tags map[string]string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		scope.SetLevel(level)
		sentry.CaptureMessage(message)
	})
}

// AddBreadcrumb adds a breadcrumb to the current scope
func AddBreadcrumb(category, message string, level sentry.Level, data map[string]interface{}) {
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: category,
		Message:  message,
		Level:    level,
		Data:     data,
		Type:     getBreadcrumbType(category),
	})
}

// getBreadcrumbType returns the appropriate breadcrumb type based on category
func getBreadcrumbType(category string) string {
	switch {
	case strings.HasPrefix(category, "http"):
		return BreadcrumbTypeHTTP
	case strings.HasPrefix(category, "error"):
		return BreadcrumbTypeError
	case strings.HasPrefix(category, "debug"):
		return BreadcrumbTypeDebug
	case strings.HasPrefix(category, "info"):
		return BreadcrumbTypeInfo
	default:
		return BreadcrumbTypeDefault
	}
}

// sanitizeEvent removes sensitive data from events
func sanitizeEvent(event *sentry.Event) *sentry.Event {
	if event == nil {
		return nil
	}

	// Sanitize request data
	if event.Request != nil {
		event.Request = sanitizeRequest(event.Request)
	}

	// Sanitize contexts
	if event.Contexts != nil {
		newContexts := make(map[string]sentry.Context)
		for key, context := range event.Contexts {
			// sentry.Context is already a map[string]interface{}, so sanitize it directly
			newContexts[key] = sanitizeMap(context)
		}
		event.Contexts = newContexts
	}

	// Sanitize extra data
	if event.Extra != nil {
		event.Extra = sanitizeMap(event.Extra)
	}

	// Sanitize tags
	if event.Tags != nil {
		event.Tags = sanitizeTags(event.Tags)
	}

	return event
}

// sanitizeRequest removes sensitive data from HTTP requests
func sanitizeRequest(request *sentry.Request) *sentry.Request {
	if request == nil {
		return nil
	}

	// Sanitize headers
	if request.Headers != nil {
		sanitized := make(map[string]string)
		for key, value := range request.Headers {
			if isSensitiveField(key) {
				sanitized[key] = "[REDACTED]"
			} else {
				sanitized[key] = value
			}
		}
		request.Headers = sanitized
	}

	// Sanitize cookies
	if request.Cookies != "" {
		request.Cookies = "[REDACTED]"
	}

	// Sanitize query string
	if request.QueryString != "" && containsSensitiveData(request.QueryString) {
		request.QueryString = sanitizeQueryString(request.QueryString)
	}

	// Sanitize body data
	if request.Data != "" && containsSensitiveData(request.Data) {
		request.Data = "[REDACTED]"
	}

	return request
}

// sanitizeMap removes sensitive data from a map
func sanitizeMap(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for key, value := range data {
		if isSensitiveField(key) {
			sanitized[key] = "[REDACTED]"
		} else {
			switch v := value.(type) {
			case string:
				if containsSensitiveData(v) {
					sanitized[key] = "[REDACTED]"
				} else {
					sanitized[key] = v
				}
			case map[string]interface{}:
				sanitized[key] = sanitizeMap(v)
			default:
				sanitized[key] = v
			}
		}
	}
	return sanitized
}

// sanitizeTags removes sensitive data from tags
func sanitizeTags(tags map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for key, value := range tags {
		if isSensitiveField(key) || containsSensitiveData(value) {
			sanitized[key] = "[REDACTED]"
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}

// sanitizeBreadcrumb removes sensitive data from breadcrumbs
func sanitizeBreadcrumb(breadcrumb *sentry.Breadcrumb) *sentry.Breadcrumb {
	if breadcrumb == nil {
		return nil
	}

	if breadcrumb.Data != nil {
		breadcrumb.Data = sanitizeMap(breadcrumb.Data)
	}

	if containsSensitiveData(breadcrumb.Message) {
		breadcrumb.Message = "[REDACTED]"
	}

	return breadcrumb
}

// sanitizeQueryString removes sensitive parameters from query strings
func sanitizeQueryString(query string) string {
	parts := strings.Split(query, "&")
	sanitized := make([]string, 0, len(parts))

	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 && isSensitiveField(keyValue[0]) {
			sanitized = append(sanitized, keyValue[0]+"=[REDACTED]")
		} else {
			sanitized = append(sanitized, part)
		}
	}

	return strings.Join(sanitized, "&")
}

// isSensitiveField checks if a field name indicates sensitive data
func isSensitiveField(field string) bool {
	fieldLower := strings.ToLower(field)
	for _, pattern := range SensitiveFieldPatterns {
		if strings.Contains(fieldLower, pattern) {
			return true
		}
	}
	return false
}

// containsSensitiveData checks if a string might contain sensitive data
func containsSensitiveData(data string) bool {
	dataLower := strings.ToLower(data)
	for _, pattern := range SensitiveFieldPatterns {
		if strings.Contains(dataLower, pattern) {
			return true
		}
	}
	return false
}

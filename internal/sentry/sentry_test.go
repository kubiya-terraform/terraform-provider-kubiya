package sentry

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "production environment",
			envValue: "production",
			expected: "production",
		},
		{
			name:     "staging environment",
			envValue: "staging",
			expected: "staging",
		},
		{
			name:     "empty environment defaults to production",
			envValue: "",
			expected: "production",
		},
		{
			name:     "invalid environment defaults to production",
			envValue: "development",
			expected: "production",
		},
		{
			name:     "uppercase staging",
			envValue: "STAGING",
			expected: "staging",
		},
		{
			name:     "whitespace staging",
			envValue: " staging ",
			expected: "staging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("KUBIYA_ENV", tt.envValue)
				defer os.Unsetenv("KUBIYA_ENV")
			} else {
				os.Unsetenv("KUBIYA_ENV")
			}

			// Test
			result := getEnvironment()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name               string
		envValue           string
		expectedEnv        string
		expectedSampleRate float64
		expectedTraceRate  float64
	}{
		{
			name:               "production config",
			envValue:           "production",
			expectedEnv:        "production",
			expectedSampleRate: ProductionErrorSampleRate,
			expectedTraceRate:  ProductionTracesSampleRate,
		},
		{
			name:               "staging config",
			envValue:           "staging",
			expectedEnv:        "staging",
			expectedSampleRate: StagingErrorSampleRate,
			expectedTraceRate:  StagingTracesSampleRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("KUBIYA_ENV", tt.envValue)
			defer os.Unsetenv("KUBIYA_ENV")

			// Test
			config := getConfig("1.0.0")

			assert.Equal(t, tt.expectedEnv, config.Environment)
			assert.Equal(t, tt.expectedSampleRate, config.SampleRate)
			assert.Equal(t, tt.expectedTraceRate, config.TracesSampleRate)
			assert.Equal(t, DSN, config.DSN)
			assert.Equal(t, EnableTracing, config.EnableTracing)
			assert.Equal(t, "terraform-provider-kubiya@1.0.0", config.Release)
		})
	}
}

func TestSanitizeMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "sanitize password field",
			input: map[string]interface{}{
				"username": "user",
				"password": "secret123",
			},
			expected: map[string]interface{}{
				"username": "user",
				"password": "[REDACTED]",
			},
		},
		{
			name: "sanitize api_key field",
			input: map[string]interface{}{
				"endpoint": "https://api.example.com",
				"api_key":  "sk-1234567890",
			},
			expected: map[string]interface{}{
				"endpoint": "https://api.example.com",
				"api_key":  "[REDACTED]",
			},
		},
		{
			name: "sanitize nested map",
			input: map[string]interface{}{
				"config": map[string]interface{}{
					"host":   "localhost",
					"secret": "mysecret",
				},
			},
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"host":   "localhost",
					"secret": "[REDACTED]",
				},
			},
		},
		{
			name: "no sensitive data",
			input: map[string]interface{}{
				"name":  "test",
				"value": "data",
			},
			expected: map[string]interface{}{
				"name":  "test",
				"value": "data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSensitiveField(t *testing.T) {
	tests := []struct {
		field     string
		sensitive bool
	}{
		{"password", true},
		{"Password", true},
		{"PASSWORD", true},
		{"api_key", true},
		{"apiKey", true},
		{"secret_token", true},
		{"authorization", true},
		{"access_token", true},
		{"refresh_token", true},
		{"private_key", true},
		{"client_secret", true},
		{"username", false},
		{"email", false},
		{"name", false},
		{"value", false},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			result := isSensitiveField(tt.field)
			assert.Equal(t, tt.sensitive, result)
		})
	}
}

func TestSanitizeQueryString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sanitize api_key in query",
			input:    "endpoint=test&api_key=secret123&format=json",
			expected: "endpoint=test&api_key=[REDACTED]&format=json",
		},
		{
			name:     "sanitize token in query",
			input:    "user=john&token=abc123",
			expected: "user=john&token=[REDACTED]",
		},
		{
			name:     "no sensitive params",
			input:    "page=1&limit=10&sort=asc",
			expected: "page=1&limit=10&sort=asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeQueryString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCaptureError(t *testing.T) {
	// This test would require mocking Sentry, which is complex
	// For now, we just ensure it doesn't panic
	ctx := context.Background()
	err := errors.New("test error")
	tags := map[string]string{
		"test": "value",
	}

	// Should not panic
	CaptureError(err, ctx, tags)

	// Test with nil error (should not panic)
	CaptureError(nil, ctx, tags)
}

func TestAddBreadcrumb(t *testing.T) {
	// This test ensures the function doesn't panic
	AddBreadcrumb("test", "test message", sentry.LevelInfo, map[string]interface{}{
		"key": "value",
	})
}

func TestGetBreadcrumbType(t *testing.T) {
	tests := []struct {
		category string
		expected string
	}{
		{"http.request", BreadcrumbTypeHTTP},
		{"error.validation", BreadcrumbTypeError},
		{"debug.trace", BreadcrumbTypeDebug},
		{"info.log", BreadcrumbTypeInfo},
		{"custom", BreadcrumbTypeDefault},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := getBreadcrumbType(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTraceResourceOperation(t *testing.T) {
	ctx := context.Background()

	// Test creating a trace
	newCtx, span := TraceResourceOperation(ctx, "kubiya_agent", "agent-123", OpResourceCreate)
	assert.NotNil(t, span)
	assert.NotEqual(t, ctx, newCtx)

	// Finish the span
	FinishSpan(span)
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()

	// No trace ID in empty context
	traceID := GetTraceID(ctx)
	assert.Empty(t, traceID)

	// With trace ID in context
	ctx = context.WithValue(ctx, contextKeyTraceID, "trace-123")
	traceID = GetTraceID(ctx)
	assert.Equal(t, "trace-123", traceID)
}

package sentry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// Context keys for storing trace data
	contextKeyTransaction contextKey = "sentry_transaction"
	contextKeySpan        contextKey = "sentry_span"
	contextKeyTraceID     contextKey = "trace_id"
)

// StartTransaction starts a new Sentry transaction
func StartTransaction(ctx context.Context, name, operation string) (context.Context, *sentry.Span) {
	// Create transaction options
	options := []sentry.SpanOption{
		sentry.WithOpName(operation),
		sentry.WithTransactionSource(sentry.SourceCustom),
	}

	// Start the transaction
	span := sentry.StartSpan(ctx, operation, options...)
	span.Description = name

	// Store transaction in context
	ctx = context.WithValue(ctx, contextKeyTransaction, span)
	ctx = context.WithValue(ctx, contextKeyTraceID, span.TraceID.String())

	return ctx, span
}

// StartSpan starts a new span within the current transaction
func StartSpan(ctx context.Context, operation, description string) (context.Context, *sentry.Span) {
	// Try to get parent span from context
	if parentSpan, ok := ctx.Value(contextKeySpan).(*sentry.Span); ok && parentSpan != nil {
		span := parentSpan.StartChild(operation)
		span.Description = description
		ctx = context.WithValue(ctx, contextKeySpan, span)
		return ctx, span
	}

	// Try to get transaction from context
	if transaction, ok := ctx.Value(contextKeyTransaction).(*sentry.Span); ok && transaction != nil {
		span := transaction.StartChild(operation)
		span.Description = description
		ctx = context.WithValue(ctx, contextKeySpan, span)
		return ctx, span
	}

	// No parent found, start a new transaction
	return StartTransaction(ctx, description, operation)
}

// FinishSpan finishes the current span
func FinishSpan(span *sentry.Span) {
	if span != nil {
		span.Finish()
	}
}

// SetSpanStatus sets the status of a span
func SetSpanStatus(span *sentry.Span, status sentry.SpanStatus) {
	if span != nil {
		span.Status = status
	}
}

// SetSpanData adds data to a span
func SetSpanData(span *sentry.Span, key string, value interface{}) {
	if span != nil {
		span.SetData(key, value)
	}
}

// SetSpanTag adds a tag to a span
func SetSpanTag(span *sentry.Span, key, value string) {
	if span != nil {
		span.SetTag(key, value)
	}
}

// GetTraceID retrieves the trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(contextKeyTraceID).(string); ok {
		return traceID
	}
	return ""
}

// GetCurrentSpan retrieves the current span from context
func GetCurrentSpan(ctx context.Context) *sentry.Span {
	if span, ok := ctx.Value(contextKeySpan).(*sentry.Span); ok {
		return span
	}
	if transaction, ok := ctx.Value(contextKeyTransaction).(*sentry.Span); ok {
		return transaction
	}
	return nil
}

// SpanFromContext is an alias for GetCurrentSpan for consistency
func SpanFromContext(ctx context.Context) *sentry.Span {
	return GetCurrentSpan(ctx)
}

// HTTPTransport is a custom HTTP transport that adds Sentry tracing headers
type HTTPTransport struct {
	Transport http.RoundTripper
}

// NewHTTPTransport creates a new HTTP transport with Sentry tracing
func NewHTTPTransport(transport http.RoundTripper) *HTTPTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &HTTPTransport{
		Transport: transport,
	}
}

// RoundTrip implements the http.RoundTripper interface
func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Start a span for the HTTP request
	ctx, span := StartSpan(ctx, OpAPICall, fmt.Sprintf("%s %s", req.Method, req.URL.Path))
	defer FinishSpan(span)

	// Add span data
	SetSpanData(span, "http.method", req.Method)
	SetSpanData(span, "http.url", req.URL.String())
	SetSpanData(span, "http.host", req.URL.Host)

	// Add trace headers for distributed tracing
	if span != nil {
		req.Header.Set(TraceHeader, span.ToSentryTrace())

		// Add baggage header if available
		if baggage := span.ToBaggage(); baggage != "" {
			req.Header.Set(BaggageHeader, baggage)
		}
	}

	// Record the start time
	startTime := time.Now()

	// Make the request
	resp, err := t.Transport.RoundTrip(req)

	// Record duration
	duration := time.Since(startTime)
	SetSpanData(span, "http.duration_ms", duration.Milliseconds())

	if err != nil {
		// Set error status
		SetSpanStatus(span, sentry.SpanStatusInternalError)
		SetSpanData(span, "error", err.Error())

		// Capture the error
		CaptureError(err, ctx, map[string]string{
			TagHTTPMethod: req.Method,
			TagHTTPURL:    req.URL.String(),
			TagErrorType:  "http_error",
		})

		return nil, err
	}

	// Set response data
	SetSpanData(span, "http.status_code", resp.StatusCode)

	// Set span status based on HTTP status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		SetSpanStatus(span, sentry.SpanStatusOK)
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
	} else if resp.StatusCode >= 500 {
		SetSpanStatus(span, sentry.SpanStatusInternalError)
	}

	return resp, nil
}

// WrapHandlerFunc wraps an HTTP handler function with Sentry tracing
func WrapHandlerFunc(handler http.HandlerFunc, operation string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract trace context from headers if present
		var spanOptions []sentry.SpanOption
		if traceHeader := r.Header.Get(TraceHeader); traceHeader != "" {
			// Continue trace from header - create span options
			spanOptions = append(spanOptions, sentry.ContinueFromHeaders(traceHeader, r.Header.Get(BaggageHeader)))
		}

		// Start transaction
		ctx, transaction := StartTransaction(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path), operation)
		defer FinishSpan(transaction)

		// Add transaction data
		SetSpanData(transaction, "http.method", r.Method)
		SetSpanData(transaction, "http.url", r.URL.String())
		SetSpanData(transaction, "http.remote_addr", r.RemoteAddr)

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		// Update request context
		r = r.WithContext(ctx)

		// Call the handler
		handler(wrapped, r)

		// Set final transaction data
		SetSpanData(transaction, "http.status_code", wrapped.statusCode)

		// Set transaction status based on HTTP status
		if wrapped.statusCode >= 200 && wrapped.statusCode < 300 {
			SetSpanStatus(transaction, sentry.SpanStatusOK)
		} else if wrapped.statusCode >= 400 && wrapped.statusCode < 500 {
			SetSpanStatus(transaction, sentry.SpanStatusInvalidArgument)
		} else if wrapped.statusCode >= 500 {
			SetSpanStatus(transaction, sentry.SpanStatusInternalError)
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// TraceResourceOperation creates a span for a resource operation
func TraceResourceOperation(ctx context.Context, resourceType, resourceID, operation string) (context.Context, *sentry.Span) {
	ctx, span := StartSpan(ctx, operation, fmt.Sprintf("%s %s", operation, resourceType))

	// Add resource-specific tags
	SetSpanTag(span, TagResourceType, resourceType)
	if resourceID != "" {
		SetSpanTag(span, TagResourceID, resourceID)
	}
	SetSpanTag(span, TagOperation, operation)

	return ctx, span
}

// TraceAPICall creates a span for an API call
func TraceAPICall(ctx context.Context, method, url string) (context.Context, *sentry.Span) {
	ctx, span := StartSpan(ctx, OpAPICall, fmt.Sprintf("%s %s", method, url))

	SetSpanData(span, "http.method", method)
	SetSpanData(span, "http.url", url)

	return ctx, span
}

// TraceValidation creates a span for validation operations
func TraceValidation(ctx context.Context, resourceType string) (context.Context, *sentry.Span) {
	return StartSpan(ctx, OpValidation, fmt.Sprintf("Validate %s", resourceType))
}

// TraceStateManagement creates a span for state management operations
func TraceStateManagement(ctx context.Context, operation string) (context.Context, *sentry.Span) {
	return StartSpan(ctx, OpStateManage, fmt.Sprintf("State: %s", operation))
}

// RecordError records an error in the current span
func RecordError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	if span := GetCurrentSpan(ctx); span != nil {
		span.SetData("error", err.Error())
		span.Status = sentry.SpanStatusInternalError
	}
}

// RecordRetry records a retry attempt in the current span
func RecordRetry(ctx context.Context, attempt int, err error) {
	if span := GetCurrentSpan(ctx); span != nil {
		span.SetTag(TagRetryCount, fmt.Sprintf("%d", attempt))
		if err != nil {
			span.SetData(fmt.Sprintf("retry.%d.error", attempt), err.Error())
		}
	}
}

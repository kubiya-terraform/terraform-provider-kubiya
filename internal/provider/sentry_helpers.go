package provider

import (
	"context"
	"fmt"

	kubiyasentry "terraform-provider-kubiya/internal/sentry"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// WithResourceTracing wraps a resource operation with Sentry tracing
func WithResourceTracing(ctx context.Context, resourceType, resourceID, operation string, fn func() error) error {
	// Start a span for the resource operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, resourceType, resourceID, operation)
	defer kubiyasentry.FinishSpan(span)

	// Add breadcrumb for the operation
	kubiyasentry.AddBreadcrumb(
		"resource.operation",
		fmt.Sprintf("%s %s %s", operation, resourceType, resourceID),
		sentry.LevelInfo,
		map[string]interface{}{
			"resource_type": resourceType,
			"resource_id":   resourceID,
			"operation":     operation,
		},
	)

	// Execute the operation
	err := fn()
	if err != nil {
		// Record the error in the span
		kubiyasentry.RecordError(ctx, err)

		// Capture the error with context
		kubiyasentry.CaptureError(err, ctx, map[string]string{
			kubiyasentry.TagResourceType: resourceType,
			kubiyasentry.TagResourceID:   resourceID,
			kubiyasentry.TagOperation:    operation,
			kubiyasentry.TagErrorType:    "resource_error",
		})
	} else {
		// Set success status
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	}

	return err
}

// TraceResourceCreate traces a resource create operation
func TraceResourceCreate(ctx context.Context, resourceType string, fn func() error) error {
	return WithResourceTracing(ctx, resourceType, "", kubiyasentry.OpResourceCreate, fn)
}

// TraceResourceRead traces a resource read operation
func TraceResourceRead(ctx context.Context, resourceType, resourceID string, fn func() error) error {
	return WithResourceTracing(ctx, resourceType, resourceID, kubiyasentry.OpResourceRead, fn)
}

// TraceResourceUpdate traces a resource update operation
func TraceResourceUpdate(ctx context.Context, resourceType, resourceID string, fn func() error) error {
	return WithResourceTracing(ctx, resourceType, resourceID, kubiyasentry.OpResourceUpdate, fn)
}

// TraceResourceDelete traces a resource delete operation
func TraceResourceDelete(ctx context.Context, resourceType, resourceID string, fn func() error) error {
	return WithResourceTracing(ctx, resourceType, resourceID, kubiyasentry.OpResourceDelete, fn)
}

// CaptureResourceError captures a resource-related error with appropriate context
func CaptureResourceError(ctx context.Context, resourceType, resourceID, operation string, err error) {
	if err == nil {
		return
	}

	kubiyasentry.CaptureError(err, ctx, map[string]string{
		kubiyasentry.TagResourceType: resourceType,
		kubiyasentry.TagResourceID:   resourceID,
		kubiyasentry.TagOperation:    operation,
		kubiyasentry.TagErrorType:    "resource_error",
	})
}

// AddResourceBreadcrumb adds a breadcrumb for resource operations
func AddResourceBreadcrumb(resourceType, resourceID, operation, message string) {
	kubiyasentry.AddBreadcrumb(
		"resource",
		message,
		sentry.LevelInfo,
		map[string]interface{}{
			"resource_type": resourceType,
			"resource_id":   resourceID,
			"operation":     operation,
		},
	)
}

// ExtractResourceInfo extracts resource type and ID from a resource request
func ExtractResourceInfo(req interface{}) (resourceType, resourceID string) {
	// This is a placeholder - actual implementation would depend on the request structure
	// You might need to adjust this based on your specific resource models
	resourceType = "unknown"
	resourceID = ""

	// Try to extract from different request types
	switch r := req.(type) {
	case *resource.ReadRequest:
		// Extract from state if available
		resourceType = "resource"
	case *resource.CreateRequest:
		resourceType = "resource"
	case *resource.UpdateRequest:
		resourceType = "resource"
	case *resource.DeleteRequest:
		resourceType = "resource"
	default:
		_ = r // Avoid unused variable warning
	}

	return resourceType, resourceID
}

// StartTerraformTransaction starts a transaction for Terraform operations
func StartTerraformTransaction(ctx context.Context, operation string) (context.Context, *sentry.Span) {
	var op string
	switch operation {
	case "apply":
		op = kubiyasentry.OpTerraformApply
	case "plan":
		op = kubiyasentry.OpTerraformPlan
	case "destroy":
		op = kubiyasentry.OpTerraformDestroy
	case "refresh":
		op = kubiyasentry.OpTerraformRefresh
	case "import":
		op = kubiyasentry.OpTerraformImport
	default:
		op = fmt.Sprintf("terraform.%s", operation)
	}

	return kubiyasentry.StartTransaction(ctx, fmt.Sprintf("Terraform %s", operation), op)
}

package vendors

import (
	"io"

	"terraform-provider-kubiya/internal/entities"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// VendorClient defines the interface that all vendor implementations must satisfy
type VendorClient interface {
	// GetVendorName returns the vendor identifier (e.g., "slack", "confluence")
	GetVendorName() string

	// PrepareCreateRequest converts the Terraform config to vendor-specific create request
	PrepareCreateRequest(config types.Map) (interface{}, error)

	// PrepareUpdateRequest converts the Terraform config to vendor-specific update request
	PrepareUpdateRequest(config types.Map) (interface{}, error)

	// ParseCreateResponse parses the vendor-specific create response
	ParseCreateResponse(resp io.Reader) (*entities.ExternalKnowledgeModel, error)

	// ParseReadResponse parses the vendor-specific read response
	ParseReadResponse(resp io.Reader) (*entities.ExternalKnowledgeModel, error)

	// ParseUpdateResponse parses the vendor-specific update response
	ParseUpdateResponse(resp io.Reader, currentModel *entities.ExternalKnowledgeModel) error

	// ParseListResponse parses the vendor-specific list response
	ParseListResponse(resp io.Reader) ([]*entities.ExternalKnowledgeModel, error)

	// ValidateConfig validates the configuration for this vendor
	ValidateConfig(config types.Map) error
}

// Registry holds all registered vendor implementations
type Registry struct {
	vendors map[string]VendorClient
}

// NewRegistry creates a new vendor registry
func NewRegistry() *Registry {
	return &Registry{
		vendors: make(map[string]VendorClient),
	}
}

// Register adds a vendor implementation to the registry
func (r *Registry) Register(vendor VendorClient) {
	r.vendors[vendor.GetVendorName()] = vendor
}

// Get retrieves a vendor implementation by name
func (r *Registry) Get(vendorName string) (VendorClient, bool) {
	vendor, ok := r.vendors[vendorName]
	return vendor, ok
}

// InitializeRegistry creates and populates the vendor registry
func InitializeRegistry() *Registry {
	registry := NewRegistry()

	// Register all vendor implementations
	registry.Register(NewSlackVendor())
	// Future vendors must be explicitly registered here:
	// registry.Register(NewConfluenceVendor())
	// registry.Register(NewNotionVendor())
	// registry.Register(NewGithubVendor())

	return registry
}

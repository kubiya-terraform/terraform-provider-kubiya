package clients

import (
	"context"
	"fmt"

	"terraform-provider-kubiya/internal/clients/vendors"
	"terraform-provider-kubiya/internal/entities"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// vendorRegistry holds all vendor implementations
var vendorRegistry = vendors.InitializeRegistry()

// extractConfigMap extracts the map from the dynamic config value
func extractConfigMap(config types.Dynamic) (types.Map, error) {
	if config.IsNull() || config.IsUnknown() {
		return types.MapNull(types.DynamicType), nil
	}

	// The underlying value could be either a map or an object
	underlyingValue := config.UnderlyingValue()

	switch v := underlyingValue.(type) {
	case types.Map:
		return v, nil
	case types.Object:
		// Convert object to map, wrapping values in Dynamic
		elements := make(map[string]attr.Value)
		for key, val := range v.Attributes() {
			// Wrap each value in Dynamic to ensure compatibility
			elements[key] = types.DynamicValue(val)
		}
		mapVal, diags := types.MapValue(types.DynamicType, elements)
		if diags.HasError() {
			return types.MapNull(types.DynamicType), fmt.Errorf("failed to convert object to map: %v", diags)
		}
		return mapVal, nil
	default:
		return types.MapNull(types.DynamicType), fmt.Errorf("config must be a map or object, got %T", underlyingValue)
	}
}

func (c *Client) ReadExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) error {
	if e == nil {
		return fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()
	path := format("/api/v1/rag/integration/%s/%s", vendor, id)

	resp, err := c.readWithJson(ctx, c.uri(path))
	if err != nil {
		return err
	}

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		return fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
	}

	model, err := vendorClient.ParseReadResponse(resp)
	if err != nil {
		return err
	}

	*e = *model
	return nil
}

func (c *Client) DeleteExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) error {
	if e == nil {
		return fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()
	path := format("/api/v1/rag/integration/%s/%s", vendor, id)

	_, err := c.deleteWithJson(ctx, c.uri(path))
	return err
}

func (c *Client) UpdateExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) error {
	if e == nil {
		return fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()
	uri := c.uri(format("/api/v1/rag/integration/%s/%s", vendor, id))

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		return fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
	}

	// Extract the map from the dynamic config
	configMap, err := extractConfigMap(e.Config)
	if err != nil {
		return err
	}

	// Prepare vendor-specific request
	requestBody, err := vendorClient.PrepareUpdateRequest(configMap)
	if err != nil {
		return err
	}

	body, err := toJson(requestBody)
	if err != nil {
		return err
	}

	resp, err := c.updateWithJson(ctx, uri, body)
	if err != nil {
		return err
	}

	// Parse vendor-specific response
	return vendorClient.ParseUpdateResponse(resp, e)
}

func (c *Client) CreateExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) (*entities.ExternalKnowledgeModel, error) {
	if e == nil {
		return nil, fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
	}

	vendor := e.Vendor.ValueString()
	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		return nil, fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
	}

	// Extract the map from the dynamic config
	configMap, err := extractConfigMap(e.Config)
	if err != nil {
		return nil, err
	}

	// Prepare vendor-specific request
	requestBody, err := vendorClient.PrepareCreateRequest(configMap)
	if err != nil {
		return nil, err
	}

	body, err := toJson(requestBody)
	if err != nil {
		return nil, err
	}

	uri := c.uri(format("/api/v1/rag/integration/%s", vendor))

	resp, err := c.createWithJson(ctx, uri, body)
	if err != nil {
		return nil, err
	}

	// Parse vendor-specific response
	return vendorClient.ParseCreateResponse(resp)
}

func (c *Client) ListExternalKnowledge(ctx context.Context, vendor string) ([]*entities.ExternalKnowledgeModel, error) {
	path := format("/api/v1/rag/integration/%s", vendor)
	uri := c.uri(path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		return nil, fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
	}

	return vendorClient.ParseListResponse(resp)
}

// This is called by the state() method
func (c *Client) externalKnowledge() ([]*vendors.BaseExternalKnowledge, error) {
	// Since we don't have a generic list endpoint, return empty list
	// Individual vendor lists should be retrieved using ListExternalKnowledge
	return []*vendors.BaseExternalKnowledge{}, nil
}

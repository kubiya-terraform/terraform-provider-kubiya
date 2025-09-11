package clients

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"

	"terraform-provider-kubiya/internal/clients/vendors"
	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"

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
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if e == nil {
		err := fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
		logger.Error("External knowledge entity is nil", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()

	logger.Debug("Reading external knowledge", map[string]interface{}{
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	// Add breadcrumb for operation tracking
	kubiyasentry.AddBreadcrumb("client", "Reading external knowledge", sentry.LevelDebug, map[string]interface{}{
		"method":                "ReadExternalKnowledge",
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	path := format("/api/v1/rag/integration/%s/%s", vendor, id)

	resp, err := c.readWithJson(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to read external knowledge", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		err := fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
		logger.Error("Unsupported vendor", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	model, err := vendorClient.ParseReadResponse(resp)
	if err != nil {
		logger.Error("Failed to parse external knowledge response", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	*e = *model

	logger.Debug("Successfully read external knowledge", map[string]interface{}{
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	return nil
}

func (c *Client) DeleteExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if e == nil {
		err := fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
		logger.Error("External knowledge entity is nil", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()

	logger.Info("Starting external knowledge deletion", map[string]interface{}{
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	// Add breadcrumb for operation tracking
	kubiyasentry.AddBreadcrumb("client", "Deleting external knowledge", sentry.LevelInfo, map[string]interface{}{
		"method":                "DeleteExternalKnowledge",
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	path := format("/api/v1/rag/integration/%s/%s", vendor, id)

	_, err := c.deleteWithJson(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to delete external knowledge", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
	} else {
		logger.Info("Successfully deleted external knowledge", map[string]interface{}{
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
	}
	return err
}

func (c *Client) UpdateExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if e == nil {
		err := fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
		logger.Error("External knowledge entity is nil", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	id := e.Id.ValueString()
	vendor := e.Vendor.ValueString()

	logger.Info("Starting external knowledge update", map[string]interface{}{
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	// Add breadcrumb for operation tracking
	kubiyasentry.AddBreadcrumb("client", "Updating external knowledge", sentry.LevelInfo, map[string]interface{}{
		"method":                "UpdateExternalKnowledge",
		"external_knowledge_id": id,
		"vendor":                vendor,
	})

	uri := c.uri(format("/api/v1/rag/integration/%s/%s", vendor, id))

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		err := fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
		logger.Error("Unsupported vendor", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	// Extract the map from the dynamic config
	configMap, err := extractConfigMap(e.Config)
	if err != nil {
		logger.Error("Failed to extract config map", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	// Prepare vendor-specific request
	requestBody, err := vendorClient.PrepareUpdateRequest(configMap)
	if err != nil {
		logger.Error("Failed to prepare vendor-specific update request", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	body, err := toJson(requestBody)
	if err != nil {
		logger.Error("Failed to marshal external knowledge data", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	resp, err := c.updateWithJson(ctx, uri, body)
	if err != nil {
		logger.Error("Failed to update external knowledge", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
		return err
	}

	// Parse vendor-specific response
	err = vendorClient.ParseUpdateResponse(resp, e)
	if err != nil {
		logger.Error("Failed to parse external knowledge update response", map[string]interface{}{
			"error":                 err.Error(),
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
	} else {
		logger.Info("Successfully updated external knowledge", map[string]interface{}{
			"external_knowledge_id": id,
			"vendor":                vendor,
		})
	}
	return err
}

func (c *Client) CreateExternalKnowledge(ctx context.Context, e *entities.ExternalKnowledgeModel) (*entities.ExternalKnowledgeModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateExternalKnowledge")
		span.SetData("client.resource_type", "external_knowledge")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating external knowledge via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateExternalKnowledge",
	})

	if e == nil {
		err := fmt.Errorf("param entity (*entities.ExternalKnowledgeModel) is nil")
		logger.Error("External knowledge entity is nil", map[string]interface{}{
			"error": err.Error(),
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	vendor := e.Vendor.ValueString()

	logger.Info("Starting external knowledge creation", map[string]interface{}{
		"vendor": vendor,
	})

	vendorClient, ok := vendorRegistry.Get(vendor)
	if !ok {
		err := fmt.Errorf("unsupported vendor: %s. Supported vendors are: slack", vendor)
		logger.Error("Unsupported vendor", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
		return nil, err
	}

	// Extract the map from the dynamic config
	configMap, err := extractConfigMap(e.Config)
	if err != nil {
		logger.Error("Failed to extract config map", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
		return nil, err
	}

	// Prepare vendor-specific request
	requestBody, err := vendorClient.PrepareCreateRequest(configMap)
	if err != nil {
		logger.Error("Failed to prepare vendor-specific create request", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
		return nil, err
	}

	body, err := toJson(requestBody)
	if err != nil {
		logger.Error("Failed to marshal external knowledge data", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
		return nil, err
	}

	uri := c.uri(format("/api/v1/rag/integration/%s", vendor))

	resp, err := c.createWithJson(ctx, uri, body)
	if err != nil {
		logger.Error("Failed to create external knowledge", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
		return nil, err
	}

	// Parse vendor-specific response
	result, err := vendorClient.ParseCreateResponse(resp)
	if err != nil {
		logger.Error("Failed to parse external knowledge create response", map[string]interface{}{
			"error":  err.Error(),
			"vendor": vendor,
		})
	} else {
		logger.Info("Successfully created external knowledge", map[string]interface{}{
			"external_knowledge_id": result.Id.ValueString(),
			"vendor":                vendor,
		})
	}
	return result, err
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
	// Get logger
	logger := kubiyasentry.GetLogger()

	logger.Debug("Fetching external knowledge list (returns empty - no generic endpoint)")

	// Since we don't have a generic list endpoint, return empty list
	// Individual vendor lists should be retrieved using ListExternalKnowledge
	return []*vendors.BaseExternalKnowledge{}, nil
}

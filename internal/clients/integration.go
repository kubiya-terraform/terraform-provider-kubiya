package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type (
	configApi struct {
		Name           string            `json:"name"`
		IsDefault      bool              `json:"is_default"`
		VendorSpecific map[string]string `json:"vendor_specific"`
	}

	integrationApi struct {
		UUID        string      `json:"uuid"`
		Name        string      `json:"name"`
		Configs     []configApi `json:"configs"`
		AuthType    string      `json:"auth_type"`
		TaskId      string      `json:"task_id"`
		ManagedBy   string      `json:"managed_by"`
		Description string      `json:"description"`
		Type        string      `json:"integration_type"`
	}
)

func newIntegration(body io.Reader) (*integrationApi, error) {
	var result integrationApi
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func diagnosticsToErrors(items diag.Diagnostics) (err error) {
	const t = "[%s] %s. Error: %s"
	for _, diagnostic := range items {
		detail := diagnostic.Detail()
		summary := diagnostic.Summary()
		severity := diagnostic.Severity()
		err = errors.Join(err, eformat(t, severity, summary, detail))
	}
	return
}

func toIntegrationApi(e *entities.IntegrationModel) (*integrationApi, error) {
	resp := &integrationApi{
		Configs:     make([]configApi, 0),
		UUID:        e.ID.ValueString(),
		Name:        e.Name.ValueString(),
		Type:        e.Type.ValueString(),
		AuthType:    e.AuthType.ValueString(),
		Description: e.Description.ValueString(),
	}

	hasDefault := false

	for _, c := range e.Configs {
		if !hasDefault {
			hasDefault = c.IsDefault.ValueBool()
		}
		config := configApi{
			Name:           c.Name.ValueString(),
			IsDefault:      c.IsDefault.ValueBool(),
			VendorSpecific: make(map[string]string),
		}

		for key, val := range c.VendorSpecific.Elements() {
			str, ok := val.(types.String)
			if ok && !str.IsNull() && !str.IsUnknown() {
				config.VendorSpecific[key] = str.ValueString()
			}
		}

		resp.Configs = append(resp.Configs, config)
	}

	if !hasDefault || len(resp.Configs) <= 0 {
		return resp, fmt.Errorf("empty configs or integration has no default config")
	}

	return resp, nil
}

func toIntegrationModel(e *integrationApi) (*entities.IntegrationModel, error) {
	var err error
	id := uuid.NewString()
	configs := make([]entities.ConfigModel, 0)

	for _, config := range e.Configs {
		vsMap := map[string]attr.Value{}
		for k, v := range config.VendorSpecific {
			vsMap[k] = types.StringValue(v)
		}

		mapValue, diags := types.MapValue(types.StringType, vsMap)
		if er := diagnosticsToErrors(diags); er != nil {
			err = errors.Join(err, er)
		}

		configs = append(configs, entities.ConfigModel{
			Name:           types.StringValue(config.Name),
			IsDefault:      types.BoolValue(config.IsDefault),
			VendorSpecific: mapValue,
		})
	}

	return &entities.IntegrationModel{
		Configs:     configs,
		ID:          types.StringValue(id),
		Name:        types.StringValue(e.Name),
		Type:        types.StringValue(e.Type),
		AuthType:    types.StringValue(e.AuthType),
		Description: types.StringValue(e.Description),
	}, err
}

func (c *Client) DeleteIntegration(ctx context.Context, e *entities.IntegrationModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteIntegration")
		span.SetData("client.resource_type", "integration")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting integration via API", sentry.LevelInfo, map[string]interface{}{
		"method": "DeleteIntegration",
	})

	if e != nil {
		name := e.Name.ValueString()
		id := e.ID.ValueString()

		logger.Info("Starting integration deletion", map[string]interface{}{
			"integration_id":   id,
			"integration_name": name,
		})

		path := format("/api/v2/integrations/%s", name)

		_, err := c.delete(ctx, c.uri(path))
		if err != nil {
			logger.Error("Failed to delete integration", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
		} else {
			logger.Info("Successfully deleted integration", map[string]interface{}{
				"integration_id":   id,
				"integration_name": name,
			})
		}
		return err
	}

	err := fmt.Errorf("param entity (*entities.IntegrationModel) is nil")
	logger.Error("Integration entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) UpdateIntegration(ctx context.Context, e *entities.IntegrationModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "UpdateIntegration")
		span.SetData("client.resource_type", "integration")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Updating integration via API", sentry.LevelInfo, map[string]interface{}{
		"method": "UpdateIntegration",
	})

	if e != nil {
		name := e.Name.ValueString()
		id := e.ID.ValueString()

		logger.Info("Starting integration update", map[string]interface{}{
			"integration_id":   id,
			"integration_name": name,
		})

		data, err := toIntegrationApi(e)
		if err != nil {
			logger.Error("Failed to convert integration to API format", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			logger.Error("Failed to marshal integration data", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		uri := c.uri(format("/api/v2/integrations/%s", name))

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to update integration", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		var r integrationApi
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode update integration response", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		e, err = toIntegrationModel(&r)
		if err != nil {
			logger.Error("Failed to convert updated integration from API response", map[string]interface{}{
				"error":            err.Error(),
				"integration_id":   id,
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
		} else {
			logger.Info("Successfully updated integration", map[string]interface{}{
				"integration_id":   id,
				"integration_name": name,
			})
		}
		return err
	}

	err := fmt.Errorf("param entity (*entities.IntegrationModel) is nil")
	logger.Error("Integration entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) ReadIntegration(ctx context.Context, name string) (*entities.IntegrationModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadIntegration")
		span.SetData("client.resource_type", "integration")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading integration via API", sentry.LevelDebug, map[string]interface{}{
		"method": "ReadIntegration",
		"name":   name,
	})

	logger.Debug("Reading integration", map[string]interface{}{
		"integration_name": name,
	})

	resp, err := c.read(ctx, c.uri(format("/api/v2/integrations/%s", name)))
	if err != nil {
		logger.Error("Failed to read integration", map[string]interface{}{
			"error":            err.Error(),
			"integration_name": name,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	r, err := newIntegration(resp)
	if err != nil {
		logger.Error("Failed to decode integration response", map[string]interface{}{
			"error":            err.Error(),
			"integration_name": name,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	entity, err := toIntegrationModel(r)
	if err != nil || entity == nil {
		if err != nil {
			logger.Error("Failed to convert integration from API response", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}
		err = eformat("Integration %s not found", name)
		logger.Error("Integration not found", map[string]interface{}{
			"error":            err.Error(),
			"integration_name": name,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Debug("Successfully read integration", map[string]interface{}{
		"integration_id":   entity.ID.ValueString(),
		"integration_name": name,
	})

	return entity, nil
}

func (c *Client) CreateIntegration(ctx context.Context, e *entities.IntegrationModel) (*entities.IntegrationModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateIntegration")
		span.SetData("client.resource_type", "integration")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating integration via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateIntegration",
		"uri":    "/api/v2/integrations",
	})

	if e != nil {
		name := e.Name.ValueString()

		logger.Info("Starting integration creation", map[string]interface{}{
			"integration_name": name,
		})

		data, err := toIntegrationApi(e)
		if err != nil {
			logger.Error("Failed to convert integration to API format", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			logger.Error("Failed to marshal integration data", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		uri := c.uri("/api/v2/integrations")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to create integration", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		var r integrationApi
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode create integration response", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		entity, err := toIntegrationModel(&r)
		if err != nil {
			logger.Error("Failed to convert created integration from API response", map[string]interface{}{
				"error":            err.Error(),
				"integration_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		logger.Info("Successfully created integration", map[string]interface{}{
			"integration_id":   entity.ID.ValueString(),
			"integration_name": name,
		})

		return entity, err
	}

	err := fmt.Errorf("param entity (*entities.IntegrationModel) is nil")
	logger.Error("Integration entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	kubiyasentry.RecordError(ctx, err)
	return nil, err
}

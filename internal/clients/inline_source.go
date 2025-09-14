package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

func newInlineSource(e *entities.InlineSourceModel) (io.Reader, error) {
	type (
		request struct {
			Name      string         `json:"name"`
			Runner    string         `json:"runner"`
			TaskId    string         `json:"task_id"`
			ManagedBy string         `json:"managed_by"`
			Tools     interface{}    `json:"inline_tools"`
			Workflows interface{}    `json:"inline_workflows"`
			Config    map[string]any `json:"dynamic_config"`
		}
	)

	req := &request{
		TaskId:    getTaskId(),
		ManagedBy: getManagedBy(),
		Config:    make(map[string]any),
		Name:      e.Name.ValueString(),
		Runner:    e.Runner.ValueString(),
	}

	if e.Tools.ValueString() != "" && e.Tools.ValueString() != "{}" {
		body := []byte(e.Tools.ValueString())
		if err := json.Unmarshal(body, &req.Tools); err != nil {
			return nil, err
		}
	}

	if e.Config.ValueString() != "" && e.Config.ValueString() != "{}" {
		body := []byte(e.Config.ValueString())
		if err := json.Unmarshal(body, &req.Config); err != nil {
			return nil, err
		}
	}

	if e.Workflows.ValueString() != "" && e.Workflows.ValueString() != "{}" {
		body := []byte(e.Workflows.ValueString())
		if err := json.Unmarshal(body, &req.Workflows); err != nil {
			return nil, err
		}
	}

	return toJson(req)
}

func parseInlineSource(r io.Reader) (*entities.InlineSourceModel, error) {
	type response struct {
		Id             string         `json:"uuid"`
		Type           string         `json:"type"`
		Url            string         `json:"url"`
		Zip            string         `json:"zip"`
		Path           string         `json:"path"`
		Name           string         `json:"name"`
		TaskId         string         `json:"task_id"`
		ManagedBy      string         `json:"managed_by"`
		AgentsCount    int            `json:"connected_agents_count"`
		ToolsCount     int            `json:"connected_tools_count"`
		WorkflowsCount int            `json:"connected_workflows_count"`
		ErrorsCount    int            `json:"errors_count"`
		Config         map[string]any `json:"dynamic_config"`
		Runner         string         `json:"runner"`
	}

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return nil, err
	}

	configData, err := json.Marshal(resp.Config)
	if err != nil {
		return nil, err
	}

	config, err := normalizeJSON(string(configData))
	if err != nil {
		return nil, err
	}

	result := &entities.InlineSourceModel{
		Config: types.StringValue(config),
		Id:     types.StringValue(resp.Id),
		Name:   types.StringValue(resp.Name),
		Type:   types.StringValue(resp.Type),
		Runner: types.StringValue(resp.Runner),
	}

	return result, nil
}

func parseNewInlineSource(r io.Reader) (*entities.InlineSourceModel, error) {
	type response struct {
		Url       string         `json:"url"`
		Type      string         `json:"type"`
		Id        string         `json:"uuid"`
		Name      string         `json:"name"`
		Runner    string         `json:"runner"`
		TaskId    string         `json:"task_id"`
		ManagedBy string         `json:"managed_by"`
		Config    map[string]any `json:"dynamic_config"`
		Errors    []struct {
			File    string `json:"file"`
			Type    string `json:"type"`
			Error   string `json:"error"`
			Details string `json:"details"`
		} `json:"errors,omitempty"`
	}

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return nil, err
	}

	if len(resp.Errors) >= 1 {
		var err error
		const t = "file: %s, type: %s, error: %s, details: %s"
		for _, e := range resp.Errors {
			err = errors.Join(err, eformat(t, e.File, e.Type, e.Error, e.Details))
		}
		return nil, err
	}

	configData, err := json.Marshal(resp.Config)
	if err != nil {
		return nil, err
	}

	config, err := normalizeJSON(string(configData))
	if err != nil {
		return nil, err
	}

	result := &entities.InlineSourceModel{
		Config: types.StringValue(config),
		Id:     types.StringValue(resp.Id),
		Name:   types.StringValue(resp.Name),
		Type:   types.StringValue(resp.Type),
		Runner: types.StringValue(resp.Runner),
	}

	return result, nil
}

func parseInlineSourceTools(r io.Reader, e *entities.InlineSourceModel) error {
	type (
		response struct {
			Id        string           `json:"uuid"`
			Type      string           `json:"type"`
			Tools     interface{}      `json:"tools"`
			Workflows []map[string]any `json:"workflows"`
			Errors    []struct {
				File    string `json:"file"`
				Type    string `json:"type"`
				Error   string `json:"error"`
				Details string `json:"details"`
			} `json:"errors,omitempty"`
		}
	)

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return err
	}

	if len(resp.Errors) >= 1 {
		var err error
		const t = "file: %s, type: %s, error: %s, details: %s"
		for _, respError := range resp.Errors {
			err = errors.Join(err, eformat(t, respError.File,
				respError.Type, respError.Error, respError.Details))
		}
		return err
	}

	if resp.Tools != nil {
		toolsData, err := json.Marshal(resp.Tools)
		if err != nil {
			return err
		}

		normalized, err := normalizeJSON(string(toolsData))
		if err != nil {
			return err
		}

		if normalized == "[]" {
			normalized = ""
		}
		e.Tools = types.StringValue(normalized)
	}

	if len(resp.Workflows) >= 1 {
		workflowList := make([]map[string]any, 0)
		for _, workflow := range resp.Workflows {
			workflowList = append(workflowList, cleanMap(workflow))
		}

		data, err := json.Marshal(workflowList)
		if err != nil {
			return err
		}

		normalized, err := normalizeJSON(string(data))
		if err != nil {
			return err
		}

		if normalized == "[]" {
			normalized = ""
		}
		e.Workflows = types.StringValue(normalized)
	}

	return nil
}

func (c *Client) DeleteInlineSource(ctx context.Context, e *entities.InlineSourceModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteInlineSource")
		span.SetData("client.resource_type", "inline_source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting inline source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "DeleteInlineSource",
	})

	const (
		requestUri = "/api/v1/sources/%s"
	)
	id := e.Id.ValueString()
	name := e.Name.ValueString()

	logger.Info("Starting inline source deletion", map[string]interface{}{
		"inline_source_id":   id,
		"inline_source_name": name,
	})

	uri := format(requestUri, id)
	resp, err := c.deleteResp(ctx, c.uri(uri))
	if err != nil {
		logger.Error("Failed to delete inline source", map[string]interface{}{
			"error":              err.Error(),
			"inline_source_id":   id,
			"inline_source_name": name,
		})
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		err := responseBodyError(resp)
		logger.Error("Delete inline source returned error status", map[string]interface{}{
			"error":              err.Error(),
			"inline_source_id":   id,
			"inline_source_name": name,
			"status_code":        resp.StatusCode,
		})
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Info("Successfully deleted inline source", map[string]interface{}{
		"inline_source_id":   id,
		"inline_source_name": name,
	})

	return nil
}

func (c *Client) UpdateInlineSource(ctx context.Context, e *entities.InlineSourceModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "UpdateInlineSource")
		span.SetData("client.resource_type", "inline_source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Updating inline source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "UpdateInlineSource",
	})

	if e != nil {
		const (
			updateUri   = "/api/v1/sources/%s"
			metadataUri = "/api/v1/sources/%s/metadata"
		)
		id := e.Id.ValueString()
		name := e.Name.ValueString()

		logger.Info("Starting inline source update", map[string]interface{}{
			"inline_source_id":   id,
			"inline_source_name": name,
		})

		uri := c.uri(format(updateUri, id))

		body, err := newInlineSource(e)
		if err != nil {
			logger.Error("Failed to convert inline source to API format", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to update inline source", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		e, err = parseInlineSource(resp)
		if err != nil {
			logger.Error("Failed to parse update inline source response", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		uri = c.uri(format(metadataUri, id))

		resp, err = c.read(ctx, uri)
		if err != nil {
			logger.Error("Failed to read inline source metadata", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		err = parseInlineSourceTools(resp, e)
		if err != nil {
			logger.Error("Failed to parse inline source tools and workflows", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		logger.Info("Successfully updated inline source", map[string]interface{}{
			"inline_source_id":   id,
			"inline_source_name": name,
		})

		return nil
	}

	err := fmt.Errorf("param entity (*entities.InlineSourceModel) is nil")
	logger.Error("Inline source entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) ReadInlineSource(ctx context.Context, id string) (*entities.InlineSourceModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadInlineSource")
		span.SetData("client.resource_type", "inline_source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading inline source via API", sentry.LevelDebug, map[string]interface{}{
		"method": "ReadInlineSource",
		"id":     id,
	})

	logger.Debug("Reading inline source", map[string]interface{}{
		"inline_source_id": id,
	})

	const (
		readUri     = "/api/v1/sources/%s"
		metadataUri = "/api/v1/sources/%s/metadata"
	)
	uri := format(readUri, id)
	resp, err := c.read(ctx, c.uri(uri))
	if err != nil {
		logger.Error("Failed to read inline source", map[string]interface{}{
			"error":            err.Error(),
			"inline_source_id": id,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	result, err := parseInlineSource(resp)
	if err != nil {
		logger.Error("Failed to parse inline source response", map[string]interface{}{
			"error":            err.Error(),
			"inline_source_id": id,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	uri = c.uri(format(metadataUri, id))
	resp, err = c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to read inline source metadata", map[string]interface{}{
			"error":            err.Error(),
			"inline_source_id": id,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	err = parseInlineSourceTools(resp, result)
	if err != nil {
		logger.Error("Failed to parse inline source tools and workflows", map[string]interface{}{
			"error":            err.Error(),
			"inline_source_id": id,
		})
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Debug("Successfully read inline source", map[string]interface{}{
		"inline_source_id":   id,
		"inline_source_name": result.Name.ValueString(),
	})

	return result, nil
}

func (c *Client) CreateInlineSource(ctx context.Context, e *entities.InlineSourceModel) (*entities.InlineSourceModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateInlineSource")
		span.SetData("client.resource_type", "inline_source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating inline source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateInlineSource",
	})

	if e != nil {
		name := e.Name.ValueString()

		logger.Info("Starting inline source creation", map[string]interface{}{
			"inline_source_name": name,
		})

		const (
			createUri   = "/api/v1/sources"
			metadataUri = "/api/v1/sources/%s/metadata"
		)

		uri := c.uri(createUri)

		body, err := newInlineSource(e)
		if err != nil {
			logger.Error("Failed to convert inline source to API format", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to create inline source", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		result, err := parseNewInlineSource(resp)
		if err != nil {
			logger.Error("Failed to parse create inline source response", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		id := result.Id.ValueString()
		uri = c.uri(format(metadataUri, id))

		resp, err = c.read(ctx, uri, "exclude_workflows_tools=true")
		if err != nil {
			logger.Error("Failed to read inline source metadata", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		err = parseInlineSourceTools(resp, result)
		if err != nil {
			logger.Error("Failed to parse inline source tools and workflows", map[string]interface{}{
				"error":              err.Error(),
				"inline_source_id":   id,
				"inline_source_name": name,
			})
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		logger.Info("Successfully created inline source", map[string]interface{}{
			"inline_source_id":   id,
			"inline_source_name": name,
		})

		return result, nil
	}

	err := fmt.Errorf("param entity (*entities.InlineSourceModel) is nil")
	logger.Error("Inline source entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	kubiyasentry.RecordError(ctx, err)
	return nil, err
}

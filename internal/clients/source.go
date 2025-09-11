package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type source struct {
	Url           string         `json:"url"`
	Id            string         `json:"uuid"`
	Name          string         `json:"name"`
	TaskId        string         `json:"task_id"`
	ManagedBy     string         `json:"managed_by"`
	DynamicConfig map[string]any `json:"dynamic_config"`
	Runner        string         `json:"runner"`
}

func newSource(body io.Reader) (*source, error) {
	var result source

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fromSource(a *source, dynamicConfigStr types.String) (*entities.SourceModel, error) {
	result := &entities.SourceModel{
		Url:    types.StringValue(a.Url),
		Id:     types.StringValue(a.Id),
		Name:   types.StringValue(a.Name),
		Runner: types.StringValue(a.Runner),
	}

	if dynamicConfigStr.ValueString() == "" {
		if len(a.DynamicConfig) >= 1 {
			marshal, err := json.Marshal(a.DynamicConfig)
			if err != nil {
				return nil, err
			}
			result.DynamicConfig = types.StringValue(string(marshal))
		}
	} else {
		result.DynamicConfig = dynamicConfigStr
	}

	return result, nil
}

func newSources(body io.Reader) ([]*source, error) {
	var result []*source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) DeleteSource(ctx context.Context, e *entities.SourceModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteSource")
		span.SetData("client.resource_type", "source")
	}

	if e == nil {
		logger.Error("DeleteSource called with nil entity")
		err := fmt.Errorf("param entity (*entities.SourceModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	sourceId := e.Id.ValueString()
	sourceUrl := e.Url.ValueString()
	sourceRunner := e.Runner.ValueString()

	logger.Info("Starting source deletion",
		"source_id", sourceId,
		"source_url", sourceUrl,
		"runner", sourceRunner)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting source via API", sentry.LevelInfo, map[string]interface{}{
		"method":     "DeleteSource",
		"source_id":  sourceId,
		"source_url": sourceUrl,
		"runner":     sourceRunner,
	})

	path := format("/api/v1/sources/%s", sourceId)
	_, err := c.delete(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to delete source",
			"source_id", sourceId,
			"source_url", sourceUrl,
			"runner", sourceRunner,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Info("Source deleted successfully",
		"source_id", sourceId,
		"source_url", sourceUrl,
		"runner", sourceRunner)

	return nil
}

func (c *Client) ReadSource(ctx context.Context, id string) (*entities.SourceModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadSource")
		span.SetData("client.resource_type", "source")
	}

	logger.Debug("Reading source",
		"source_id", id)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading source via API", sentry.LevelInfo, map[string]interface{}{
		"method":    "ReadSource",
		"source_id": id,
	})

	path := format("/api/v1/sources/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to read source",
			"source_id", id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	result, err := newSource(resp)
	if err != nil {
		logger.Error("Failed to decode source response",
			"source_id", id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	entity, err := fromSource(result, types.StringValue("{}"))
	if err != nil {
		logger.Error("Failed to convert source",
			"source_id", id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	if result.DynamicConfig == nil {
		entity.DynamicConfig = types.StringValue("{}")
	}

	logger.Debug("Source read successfully",
		"source_id", id,
		"source_url", result.Url,
		"runner", result.Runner)

	return entity, nil
}

func (c *Client) CreateSource(ctx context.Context, e *entities.SourceModel) (*entities.SourceModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateSource")
		span.SetData("client.resource_type", "source")
	}

	if e == nil {
		logger.Error("CreateSource called with nil entity")
		err := fmt.Errorf("param entity (*entities.SourceModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	sourceUrl := e.Url.ValueString()
	sourceRunner := e.Runner.ValueString()
	dynamicConfig := e.DynamicConfig.ValueString()

	logger.Info("Starting source creation",
		"source_url", sourceUrl,
		"runner", sourceRunner,
		"has_dynamic_config", dynamicConfig != "")

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating source via API", sentry.LevelInfo, map[string]interface{}{
		"method":     "CreateSource",
		"uri":        "/api/v1/sources",
		"source_url": sourceUrl,
		"runner":     sourceRunner,
	})

	uri := c.uri("/api/v1/sources")

	data := &source{
		TaskId:        getTaskId(),
		ManagedBy:     getManagedBy(),
		Url:           sourceUrl,
		DynamicConfig: make(map[string]any),
		Runner:        sourceRunner,
	}

	if dynamicConfig != "" {
		if err := json.Unmarshal([]byte(dynamicConfig), &data.DynamicConfig); err != nil {
			logger.Error("Failed to unmarshal dynamic config",
				"source_url", sourceUrl,
				"runner", sourceRunner,
				"dynamic_config", dynamicConfig,
				"error", err.Error())
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}
	}

	body, err := toJson(data)
	if err != nil {
		logger.Error("Failed to marshal source data",
			"source_url", sourceUrl,
			"runner", sourceRunner,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	qps := []string{fmt.Sprintf("runner=%s", data.Runner)}

	resp, err := c.create(ctx, uri, body, qps...)
	if err != nil {
		logger.Error("Failed to create source via API",
			"source_url", sourceUrl,
			"runner", sourceRunner,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	result, err := newSource(resp)
	if err != nil {
		logger.Error("Failed to decode create source response",
			"source_url", sourceUrl,
			"runner", sourceRunner,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	returnSource, err := fromSource(result, e.DynamicConfig)
	if err != nil {
		logger.Error("Failed to convert created source",
			"source_url", sourceUrl,
			"runner", sourceRunner,
			"source_id", result.Id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	if result.DynamicConfig == nil {
		returnSource.DynamicConfig = types.StringValue("{}")
	}

	logger.Info("Source created successfully",
		"source_id", result.Id,
		"source_url", sourceUrl,
		"runner", sourceRunner)

	return returnSource, nil
}

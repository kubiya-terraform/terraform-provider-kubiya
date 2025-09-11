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
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteSource")
		span.SetData("client.resource_type", "source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "DeleteSource",
	})

	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/sources/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
		}
		return err
	}

	err := fmt.Errorf("param entity (*entities.SourceModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) ReadSource(ctx context.Context, id string) (*entities.SourceModel, error) {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadSource")
		span.SetData("client.resource_type", "source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "ReadSource",
	})

	path := format("/api/v1/sources/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	result, err := newSource(resp)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	entity, err := fromSource(result, types.StringValue("{}"))
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	if result.DynamicConfig == nil {
		entity.DynamicConfig = types.StringValue("{}")
	}

	return entity, nil
}

func (c *Client) CreateSource(ctx context.Context, e *entities.SourceModel) (*entities.SourceModel, error) {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateSource")
		span.SetData("client.resource_type", "source")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating source via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateSource",
		"uri":    "/api/v1/sources",
	})

	if e != nil {
		uri := c.uri("/api/v1/sources")

		data := &source{
			TaskId:        getTaskId(),
			ManagedBy:     getManagedBy(),
			Url:           e.Url.ValueString(),
			DynamicConfig: make(map[string]any),
			Runner:        e.Runner.ValueString(),
		}

		if e.DynamicConfig.ValueString() != "" {
			if err := json.Unmarshal([]byte(e.DynamicConfig.ValueString()), &data.DynamicConfig); err != nil {
				kubiyasentry.RecordError(ctx, err)
				return nil, err
			}
		}

		body, err := toJson(data)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		qps := []string{fmt.Sprintf("runner=%s", data.Runner)}

		resp, err := c.create(ctx, uri, body, qps...)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		result, err := newSource(resp)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		returnSource, err := fromSource(result, e.DynamicConfig)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		if result.DynamicConfig == nil {
			returnSource.DynamicConfig = types.StringValue("{}")
		}

		return returnSource, nil
	}

	err := fmt.Errorf("param entity (*entities.SourceModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return nil, err
}

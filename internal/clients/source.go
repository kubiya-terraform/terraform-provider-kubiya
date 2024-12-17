package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"terraform-provider-kubiya/internal/entities"
)

type source struct {
	Url           string         `json:"url"`
	Id            string         `json:"uuid"`
	Name          string         `json:"name"`
	TaskId        string         `json:"task_id"`
	ManagedBy     string         `json:"managed_by"`
	DynamicConfig map[string]any `json:"dynamic_config"`
}

func newSource(body io.Reader) (*source, error) {
	var result source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fromSource(a *source) (*entities.SourceModel, error) {
	result := &entities.SourceModel{
		Url:  types.StringValue(a.Url),
		Id:   types.StringValue(a.Id),
		Name: types.StringValue(a.Name),
	}

	dynamicConfig, err := convertToTypesMap(a.DynamicConfig)
	if err != nil {
		return nil, err
	}

	result.DynamicConfig = dynamicConfig

	return result, err
}

func newSources(body io.Reader) ([]*source, error) {
	var result []*source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) DeleteSource(ctx context.Context, e *entities.SourceModel) error {
	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/sources/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		return err
	}

	return fmt.Errorf("param entity (*entities.SourceModel) is nil")
}

func (c *Client) ReadSource(ctx context.Context, id string) (*entities.SourceModel, error) {
	path := format("/api/v1/sources/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		return nil, err
	}

	r, err := newSource(resp)
	if err != nil {
		return nil, err
	}

	return fromSource(r)
}

func (c *Client) CreateSource(ctx context.Context, e *entities.SourceModel) (*entities.SourceModel, error) {
	if e != nil {
		uri := c.uri("/api/v1/sources")

		DynamicConfig, err := processDynamicConfig(e.DynamicConfig)
		if err != nil {
			return nil, err
		}

		data := &source{
			TaskId:        getTaskId(),
			ManagedBy:     getManagedBy(),
			Url:           e.Url.ValueString(),
			DynamicConfig: DynamicConfig,
		}

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		result, err := newSource(resp)
		if err != nil {
			return nil, err
		}

		return fromSource(result)
	}

	return nil, fmt.Errorf("param entity (*entities.SourceModel) is nil")
}

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"strconv"
	"terraform-provider-kubiya/internal/entities"
)

type source struct {
	Url           string            `json:"url"`
	Id            string            `json:"uuid"`
	Name          string            `json:"name"`
	TaskId        string            `json:"task_id"`
	ManagedBy     string            `json:"managed_by"`
	DynamicConfig map[string]string `json:"dynamic_config"`
}

func newSource(body io.Reader) (*source, error) {
	var result source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println(result.DynamicConfig)
	for key, value := range result.DynamicConfig {
		cleanedString, err := strconv.Unquote(value)
		if err != nil {
			return nil, err
		}
		result.DynamicConfig[key] = cleanedString
	}

	return &result, nil
}

func fromSource(a *source) (*entities.SourceModel, error) {
	result := &entities.SourceModel{
		Url:  types.StringValue(a.Url),
		Id:   types.StringValue(a.Id),
		Name: types.StringValue(a.Name),
	}
	var err error
	result.DynamicConfig = toMapType(a.DynamicConfig, err)

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

		data := &source{
			TaskId:        getTaskId(),
			ManagedBy:     getManagedBy(),
			Url:           e.Url.ValueString(),
			DynamicConfig: make(map[string]string),
		}

		for key, value := range e.DynamicConfig.Elements() {
			data.DynamicConfig[key] = value.String()
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

package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type source struct {
	Url  string `json:"url"`
	Id   string `json:"uuid"`
	Name string `json:"name"`
}

func newSource(body io.Reader) (*source, error) {
	var result source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func toSource(e *entities.SourceModel) *source {
	result := &source{
		Id:   e.Id.ValueString(),
		Url:  e.Url.ValueString(),
		Name: e.Name.ValueString(),
	}

	return result
}

func fromSource(a *source) *entities.SourceModel {
	result := &entities.SourceModel{
		Id:   types.StringValue(a.Id),
		Url:  types.StringValue(a.Url),
		Name: types.StringValue(a.Name),
	}

	return result
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

	return fromSource(r), nil
}

func (c *Client) CreateSource(ctx context.Context, e *entities.SourceModel) (*entities.SourceModel, error) {
	if e != nil {
		data := toSource(e)

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/sources")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		result, err := newSource(resp)
		if err != nil {
			return nil, err
		}

		return fromSource(result), nil
	}

	return nil, fmt.Errorf("param entity (*entities.SourceModel) is nil")
}

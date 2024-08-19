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
	Url       string `json:"url"`
	Id        string `json:"uuid"`
	Name      string `json:"name"`
	TaskId    string `json:"task_id"`
	ManagedBy string `json:"managed_by"`
}

func newSource(body io.Reader) (*source, error) {
	var result source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func newSources(body io.Reader) ([]*source, error) {
	var result []*source
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func fromSource(a *source) *entities.SourceModel {
	result := &entities.SourceModel{
		Url:  types.StringValue(a.Url),
		Id:   types.StringValue(a.Id),
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
		uri := c.uri("/api/v1/sources")

		resp, err := c.read(ctx, uri)
		if err != nil {
			return nil, err
		}

		results, err := newSources(resp)
		if err != nil {
			return nil, err
		}

		exist := false
		srUrl := e.Url.ValueString()
		for _, sr := range results {
			if exist = sr.Url == srUrl; exist {
				break
			}
		}

		if exist {
			return nil, eformat("source already exists")
		}

		data := &source{
			TaskId:    getTaskId(),
			ManagedBy: getManagedBy(),
			Url:       e.Url.ValueString(),
		}

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		_, err = c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		resp, err = c.read(ctx, uri)
		if err != nil {
			return nil, err
		}

		results, err = newSources(resp)
		if err != nil {
			return nil, err
		}

		var result *entities.SourceModel
		for _, sr := range results {
			if sr.Url == srUrl {
				result = fromSource(sr)
				break
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("param entity (*entities.SourceModel) is nil")
}

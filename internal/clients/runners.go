package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type runner struct {
	Key     string `json:"key"`
	Url     string `json:"url"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

func (c *Client) ReadRunner(ctx context.Context, entity *entities.RunnerModel) error {
	if entity != nil {
		const (
			uri = "/api/v3/runners/%s/describe"
		)

		name := entity.Name.ValueString()

		resp, err := c.read(ctx, c.uri(format(uri, name)))
		if err != nil {
			return err
		}

		var r *runner
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return err
		}

		entity = &entities.RunnerModel{
			Key:     types.StringValue(r.Key),
			Url:     types.StringValue(r.Url),
			Name:    types.StringValue(r.Name),
			Subject: types.StringValue(r.Subject),
		}

		return err
	}

	return fmt.Errorf("param entity (*entities.RunnerModel) is nil")
}

func (c *Client) DeleteRunner(ctx context.Context, entity *entities.RunnerModel) error {
	if entity != nil {
		const (
			uri = "/api/v3/runners/%s"
		)
		name := entity.Name.ValueString()

		_, err := c.delete(ctx, c.uri(format(uri, name)))
		return err
	}

	return fmt.Errorf("param entity (*entities.RunnerModel) is nil")
}

func (c *Client) CreateRunner(ctx context.Context, entity *entities.RunnerModel) (*entities.RunnerModel, error) {
	if entity != nil {
		const (
			uri = "/api/v3/runners/%s"
		)

		path := entity.Path.ValueString()
		name := entity.Name.ValueString()

		resp, err := c.create(ctx, c.uri(format(uri, name)), nil)
		if err != nil {
			return nil, err
		}

		var r *runner
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		if len(path) >= 1 {
			path = toPathYaml(path, name)
			if err = c.downloadFile(r.Url, path); err != nil {
				return nil, err
			}

			entity.Path = types.StringValue(path)
		}

		runners, err := c.runners()
		if err != nil {
			return nil, err
		}

		for _, item := range runners {
			if equal(item.Name, name) {
				entity.Url = types.StringValue(r.Url)
				entity.Key = types.StringValue(item.Key)
				entity.Name = types.StringValue(item.Name)
				entity.Subject = types.StringValue(item.Subject)
				break
			}
		}

		return entity, err
	}

	return nil, fmt.Errorf("param entity (*entities.RunnerModel) is nil")
}

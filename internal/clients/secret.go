package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type (
	Secret struct {
		Name        string    `json:"name"`
		Value       string    `json:"value"`
		Description string    `json:description`
		CreatedAt   time.Time `json:"created_at"`
		CreatedBy   string    `json:"created_by"`
	}
)

func toSecret(s *entities.SecretModel, cs *state) *secret {
	ret := &secret{
		Name:        s.Name.ValueString(),
		Value:       s.Value.ValueString(),
		Description: s.Description.ValueString(),
		CreatedBy:   s.CreatedBy.ValueString(),
		CreatedAt:   s.CreatedAt.ValueString(),
	}

	return ret
}

func fromSecret(s *secret, cs *state) *entities.SecretModel {
	ret := &entities.SecretModel{
		CreatedAt:   types.StringValue(""),
		CreatedBy:   types.StringValue(""),
		Name:        types.StringValue(s.Name),
		Value:       types.StringValue(s.Value),
		Description: types.StringValue(s.Description),
	}

	return ret
}

func (c *Client) ReadSecret(_ context.Context, entity *entities.SecretModel) error {
	if entity != nil {
		// cs, err := c.state()
		// if err != nil {
		// 	return err
		// }
		fmt.Println(c.secrets())

		entity = &entities.SecretModel{
			Name:        types.StringValue("a name"),
			Value:       types.StringValue("a value"),
			Description: types.StringValue("a description"),
		}

		return nil
	}

	return fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

func (c *Client) DeleteSecret(ctx context.Context, entity *entities.SecretModel) error {
	if entity != nil {
		const (
			ok     = ""
			path   = "/api/v2/secrets/%s"
			errMsg = "failed to delete secret - %s"
		)

		uri := c.uri(fmt.Sprintf(path, entity.Name))
		resp, err := c.delete(ctx, uri)
		if err != nil {
			return err
		}

		r := &struct {
			Result string `json:"result"`
		}{}

		err = json.NewDecoder(resp).Decode(&r)
		if err != nil || r == nil {
			if err != nil {
				return err
			}
			return fmt.Errorf(errMsg, entity.Name)
		}

		if strings.Contains(r.Result, ok) {
			return nil
		}

		return fmt.Errorf(errMsg, entity.Name)
	}

	return fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

func (c *Client) UpdateSecret(ctx context.Context, entity *entities.SecretModel) error {
	if entity != nil {
		const (
			path = "/api/v2/secrets/%s"
		)

		cs, err := c.state()
		if err != nil {
			return err
		}

		uri := c.uri(format(path, entity.Name))

		data := toSecret(entity, cs)

		body, err := toJson(data)
		if err != nil {
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			return err
		}

		var r *secret
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return err
		}

		entity = fromSecret(r, cs)

		return err
	}
	return fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

func (c *Client) CreateSecret(ctx context.Context, entity *entities.SecretModel) (*entities.SecretModel, error) {
	if entity != nil {
 		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/secret/create_secret")

		payload := map[string]string{
			"secret_name":  entity.Name.String(),
			"secret_value": entity.Value.String(),
			"description":  entity.Description.String(),
		}
		body, err := toJson(payload)
		if err != nil {
			return nil, err
		}
		resp, err := c.create(ctx, uri, nil, body)
		if err != nil {
			return nil, err
		}

		var r *secret
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		return fromSecret(r, cs), err
	}

	return nil, fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

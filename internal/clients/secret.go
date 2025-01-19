package clients

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
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

func toSecret(s *entities.SecretModel) *secret {
	ret := &secret{
		Name:        s.Name.ValueString(),
		Value:       s.Value.ValueString(),
		Description: s.Description.ValueString(),
		CreatedBy:   s.CreatedBy.ValueString(),
		CreatedAt:   s.CreatedAt.ValueString(),
	}

	return ret
}

func fromSecret(s *secret) *entities.SecretModel {
	ret := &entities.SecretModel{
		CreatedAt:   types.StringValue(s.CreatedAt),
		CreatedBy:   types.StringValue(s.CreatedBy),
		Name:        types.StringValue(s.Name),
		Value:       types.StringValue(s.Value),
		Description: types.StringValue(s.Description),
	}

	return ret
}

func (c *Client) ReadSecret(ctx context.Context, entity *entities.SecretModel) error {
	if entity != nil {
		secretname := entity.Name.ValueString()
		if entity.Name.ValueString() == "" {
			return fmt.Errorf("secret name is empty")
		}

		// get secret metadata
		uri := c.uri(fmt.Sprintf("/api/v2/secrets/%s", secretname))

		resp, err := c.read(ctx, uri)
		if err != nil {
			return err
		}
		s := &secret{}
		err = json.NewDecoder(resp).Decode(s)
		if err != nil {
			return fmt.Errorf("failed to decode secret metadata - %s", err)
		}

		// get secret value
		uri = c.uri(fmt.Sprintf("/api/v2/secrets/get_value/%s", secretname))
		resp, err = c.read(ctx, uri)
		if err != nil {
			return err
		}
		var secretValueEncoded string
		err = json.NewDecoder(resp).Decode(&secretValueEncoded)
		if err != nil {
			return fmt.Errorf("failed to read secret value - %s", err)
		}
		secretValue, err := b64.StdEncoding.DecodeString(string(secretValueEncoded))
		if err != nil {
			return fmt.Errorf("failed to decode secret value - %s", err)
		}

		s.Value = string(secretValue)
		*entity = *fromSecret(s)

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

		uri := c.uri(fmt.Sprintf(path, entity.Name.ValueString()))
		resp, err := c.delete(ctx, uri)
		if err != nil {
			return err
		}

		r := &struct {
			Error string `json:"error"`
		}{}

		err = json.NewDecoder(resp).Decode(&r)
		if err != nil || r == nil {
			if err != nil {
				return err
			}
			return fmt.Errorf(errMsg, entity.Name)
		}

		if r.Error != "" {
			return fmt.Errorf(errMsg, entity.Name)
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

		uri := c.uri(format(path, entity.Name.ValueString()))

		data := toSecret(entity)

		body, err := toJson(data)
		if err != nil {
			return err
		}
		resp, err := c.update(ctx, uri, body)
		if err != nil {
			return err
		}
		obj := map[string]any{}
		err = json.NewDecoder(resp).Decode(&obj)
		if err != nil {
			return err
		}
		if obj["error"] != nil {
			return fmt.Errorf("failed to update secret - %s", obj["error"])
		}

		return nil
	}
	return fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

func (c *Client) CreateSecret(ctx context.Context, entity *entities.SecretModel) (*entities.SecretModel, error) {
	if entity != nil {

		uri := c.uri("/api/v2/secrets")
		payload := map[string]string{
			"name":        entity.Name.ValueString(),
			"value":       entity.Value.ValueString(),
			"description": entity.Description.ValueString(),
		}
		body, err := toJson(payload)
		if err != nil {
			return nil, err
		}
		resp, err := c.create(ctx, uri, nil, body)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, fmt.Errorf("response is nil")
		}

		return entity, nil
	}

	return nil, fmt.Errorf("param entity (*entities.SecretModel) is nil")
}

package clients

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type (
	Secret struct {
		Name        string    `json:"name"`
		Value       string    `json:"value"`
		CreatedAt   time.Time `json:"created_at"`
		CreatedBy   string    `json:"created_by"`
		Description string    `json:"description"`
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
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadSecret")
		span.SetData("client.resource_type", "secret")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading secret via API", sentry.LevelInfo, map[string]interface{}{
		"method": "ReadSecret",
	})

	if entity != nil {
		secretname := entity.Name.ValueString()
		if entity.Name.ValueString() == "" {
			err := fmt.Errorf("secret name is empty")
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		// get secret metadata
		uri := c.uri(fmt.Sprintf("/api/v2/secrets/%s", secretname))

		resp, err := c.read(ctx, uri)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		s := &secret{}
		err = json.NewDecoder(resp).Decode(s)
		if err != nil {
			err = fmt.Errorf("failed to decode secret metadata - %s", err)
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		// get secret value
		uri = c.uri(fmt.Sprintf("/api/v2/secrets/get_value/%s", secretname))
		resp, err = c.read(ctx, uri)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		var secretValueEncoded string
		err = json.NewDecoder(resp).Decode(&secretValueEncoded)
		if err != nil {
			err = fmt.Errorf("failed to read secret value - %s", err)
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		secretValue, err := b64.StdEncoding.DecodeString(string(secretValueEncoded))
		if err != nil {
			err = fmt.Errorf("failed to decode secret value - %s", err)
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		s.Value = string(secretValue)
		*entity = *fromSecret(s)

		return nil
	}

	err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) DeleteSecret(ctx context.Context, entity *entities.SecretModel) error {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteSecret")
		span.SetData("client.resource_type", "secret")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting secret via API", sentry.LevelInfo, map[string]interface{}{
		"method": "DeleteSecret",
	})

	if entity != nil {
		const (
			path   = "/api/v2/secrets/%s"
			errMsg = "failed to delete secret - %s"
		)

		uri := c.uri(fmt.Sprintf(path, entity.Name.ValueString()))
		resp, err := c.delete(ctx, uri)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		r := &struct {
			Error string `json:"error"`
		}{}

		err = json.NewDecoder(resp).Decode(&r)
		if err != nil || r == nil {
			if err != nil {
				kubiyasentry.RecordError(ctx, err)
				return err
			}
			err = fmt.Errorf(errMsg, entity.Name)
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		if r.Error != "" {
			err = fmt.Errorf(errMsg, entity.Name)
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		err = fmt.Errorf(errMsg, entity.Name)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) UpdateSecret(ctx context.Context, entity *entities.SecretModel) error {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "UpdateSecret")
		span.SetData("client.resource_type", "secret")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Updating secret via API", sentry.LevelInfo, map[string]interface{}{
		"method": "UpdateSecret",
	})

	if entity != nil {
		const (
			path = "/api/v2/secrets/%s"
		)

		uri := c.uri(format(path, entity.Name.ValueString()))

		data := toSecret(entity)

		body, err := toJson(data)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		resp, err := c.update(ctx, uri, body)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		obj := map[string]any{}
		err = json.NewDecoder(resp).Decode(&obj)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		if obj["error"] != nil {
			err = fmt.Errorf("failed to update secret - %s", obj["error"])
			kubiyasentry.RecordError(ctx, err)
			return err
		}

		return nil
	}
	err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) CreateSecret(ctx context.Context, entity *entities.SecretModel) (*entities.SecretModel, error) {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateSecret")
		span.SetData("client.resource_type", "secret")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating secret via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateSecret",
		"uri":    "/api/v2/secrets",
	})

	if entity != nil {

		uri := c.uri("/api/v2/secrets")
		payload := map[string]string{
			"name":        entity.Name.ValueString(),
			"value":       entity.Value.ValueString(),
			"description": entity.Description.ValueString(),
		}
		body, err := toJson(payload)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}
		resp, err := c.create(ctx, uri, body)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}
		if resp == nil {
			err = fmt.Errorf("response is nil")
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		return entity, nil
	}

	err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return nil, err
}

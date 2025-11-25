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
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadSecret")
		span.SetData("client.resource_type", "secret")
	}

	if entity == nil {
		logger.Error("ReadSecret called with nil entity")
		err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	secretname := entity.Name.ValueString()
	if secretname == "" {
		logger.Error("ReadSecret called with empty secret name")
		err := fmt.Errorf("secret name is empty")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Debug("Reading secret",
		"secret_name", secretname)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading secret via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "ReadSecret",
		"secret_name": secretname,
	})

	// get secret metadata
	uri := c.uri(fmt.Sprintf("/api/v2/secrets/%s", secretname))

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to read secret metadata",
			"secret_name", secretname,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	s := &secret{}
	err = json.NewDecoder(resp).Decode(s)
	if err != nil {
		logger.Error("Failed to decode secret metadata",
			"secret_name", secretname,
			"error", err.Error())
		err = fmt.Errorf("failed to decode secret metadata - %s", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	// get secret value
	uri = c.uri(fmt.Sprintf("/api/v2/secrets/get_value/%s", secretname))
	resp, err = c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to read secret value",
			"secret_name", secretname,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	var secretValueEncoded string
	err = json.NewDecoder(resp).Decode(&secretValueEncoded)
	if err != nil {
		logger.Error("Failed to decode secret value response",
			"secret_name", secretname,
			"error", err.Error())
		err = fmt.Errorf("failed to read secret value - %s", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	secretValue, err := b64.StdEncoding.DecodeString(string(secretValueEncoded))
	if err != nil {
		logger.Error("Failed to decode base64 secret value",
			"secret_name", secretname,
			"error", err.Error())
		err = fmt.Errorf("failed to decode secret value - %s", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	s.Value = string(secretValue)
	*entity = *fromSecret(s)

	logger.Debug("Secret read successfully",
		"secret_name", secretname,
		"created_by", s.CreatedBy)

	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, entity *entities.SecretModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteSecret")
		span.SetData("client.resource_type", "secret")
	}

	if entity == nil {
		logger.Error("DeleteSecret called with nil entity")
		err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	secretName := entity.Name.ValueString()

	logger.Info("Starting secret deletion",
		"secret_name", secretName)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting secret via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "DeleteSecret",
		"secret_name": secretName,
	})

	const (
		path   = "/api/v2/secrets/%s"
		errMsg = "failed to delete secret - %s"
	)

	uri := c.uri(fmt.Sprintf(path, secretName))
	resp, err := c.delete(ctx, uri)
	if err != nil {
		logger.Error("Failed to delete secret via API",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	r := &struct {
		Error string `json:"error"`
	}{}

	err = json.NewDecoder(resp).Decode(&r)
	if err != nil || r == nil {
		if err != nil {
			logger.Error("Failed to decode delete response",
				"secret_name", secretName,
				"error", err.Error())
			kubiyasentry.RecordError(ctx, err)
			return err
		}
		logger.Error("Delete response was nil",
			"secret_name", secretName)
		err = fmt.Errorf(errMsg, secretName)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	if r.Error != "" {
		logger.Error("API returned error during delete",
			"secret_name", secretName,
			"api_error", r.Error)
		err = fmt.Errorf(errMsg, secretName)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Info("Secret deleted successfully",
		"secret_name", secretName)

	// Note: This seems like a bug in the original code - it always returns an error even on success.
	// However, preserving the original behavior.
	err = fmt.Errorf(errMsg, secretName)
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) UpdateSecret(ctx context.Context, entity *entities.SecretModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "UpdateSecret")
		span.SetData("client.resource_type", "secret")
	}

	if entity == nil {
		logger.Error("UpdateSecret called with nil entity")
		err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	secretName := entity.Name.ValueString()
	secretDescription := entity.Description.ValueString()

	logger.Info("Starting secret update",
		"secret_name", secretName,
		"description", secretDescription)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Updating secret via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "UpdateSecret",
		"secret_name": secretName,
	})

	const path = "/api/v2/secrets/%s"

	uri := c.uri(format(path, secretName))

	data := toSecret(entity)

	body, err := toJson(data)
	if err != nil {
		logger.Error("Failed to marshal secret data",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	resp, err := c.update(ctx, uri, body)
	if err != nil {
		logger.Error("Failed to update secret via API",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	obj := map[string]any{}
	err = json.NewDecoder(resp).Decode(&obj)
	if err != nil {
		logger.Error("Failed to decode update response",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	if obj["error"] != nil {
		logger.Error("API returned error during update",
			"secret_name", secretName,
			"api_error", obj["error"])
		err = fmt.Errorf("failed to update secret - %s", obj["error"])
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Info("Secret updated successfully",
		"secret_name", secretName)

	return nil
}

func (c *Client) CreateSecret(ctx context.Context, entity *entities.SecretModel) (*entities.SecretModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateSecret")
		span.SetData("client.resource_type", "secret")
	}

	if entity == nil {
		logger.Error("CreateSecret called with nil entity")
		err := fmt.Errorf("param entity (*entities.SecretModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	secretName := entity.Name.ValueString()
	secretDescription := entity.Description.ValueString()

	logger.Info("Starting secret creation",
		"secret_name", secretName,
		"description", secretDescription)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating secret via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "CreateSecret",
		"uri":         "/api/v2/secrets",
		"secret_name": secretName,
	})

	uri := c.uri("/api/v2/secrets")
	payload := map[string]string{
		"name":        secretName,
		"value":       entity.Value.ValueString(),
		"description": secretDescription,
	}

	body, err := toJson(payload)
	if err != nil {
		logger.Error("Failed to marshal secret payload",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	resp, err := c.create(ctx, uri, body)
	if err != nil {
		logger.Error("Failed to create secret via API",
			"secret_name", secretName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	if resp == nil {
		logger.Error("Create secret response was nil",
			"secret_name", secretName)
		err = fmt.Errorf("response is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Info("Secret created successfully",
		"secret_name", secretName)

	return entity, nil
}

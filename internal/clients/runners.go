package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getsentry/sentry-go"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

// runner is used internally by the client
type runner struct {
	Name                string  `json:"name"`
	Subject             string  `json:"subject"`
	Namespace           string  `json:"namespace"`
	RunnerType          string  `json:"runner_type"`
	UserKeyId           string  `json:"user_key_id"`
	Version             int     `json:"version"`
	Description         string  `json:"description"`
	AuthenticationType  string  `json:"authentication_type"`
	ManagedBy           string  `json:"managed_by"`
	TaskId              string  `json:"task_id"`
	WssUrl              string  `json:"wss_url"`
	KubernetesNamespace string  `json:"kubernetes_namespace"`
	GatewayUrl          *string `json:"gateway_url"`
	GatewayPassword     *string `json:"gateway_password"`
	AgentManagerHealth  struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"agent_manager_health"`
	RunnerHealth struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"runner_health"`
	ToolManagerHealth struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"tool_manager_health"`
}

func (c *Client) ReadRunner(ctx context.Context, entity *entities.RunnerModel) error {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadRunner")
		span.SetData("client.resource_type", "runner")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading runner via API", sentry.LevelInfo, map[string]interface{}{
		"method": "ReadRunner",
		"uri":    "/api/v3/runners/%s/describe",
	})

	if entity == nil {
		err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	const uri = "/api/v3/runners/%s/describe"
	name := entity.Name.ValueString()

	// Add runner name to span
	if span != nil {
		span.SetData("runner.name", name)
	}

	reqUri := c.uri(format(uri, name))

	// Add API call details to span
	if span != nil {
		span.SetData("api.uri", reqUri)
		span.SetData("api.method", "GET")
	}

	kubiyasentry.AddBreadcrumb("client", "Making API request", sentry.LevelInfo, map[string]interface{}{
		"uri":    reqUri,
		"method": "GET",
	})

	resp, err := c.read(ctx, reqUri)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "API request failed", sentry.LevelError, map[string]interface{}{
			"uri":   reqUri,
			"error": err.Error(),
		})
		return err
	}

	var r runner
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "Failed to decode response", sentry.LevelError, map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	entity.Name = types.StringValue(r.Name)
	entity.RunnerType = types.StringValue(r.RunnerType)

	// Success breadcrumb
	kubiyasentry.AddBreadcrumb("client", "Successfully read runner", sentry.LevelInfo, map[string]interface{}{
		"runner_name": r.Name,
		"runner_type": r.RunnerType,
	})

	return nil
}

func (c *Client) DeleteRunner(ctx context.Context, entity *entities.RunnerModel) error {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteRunner")
		span.SetData("client.resource_type", "runner")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting runner via API", sentry.LevelInfo, map[string]interface{}{
		"method": "DeleteRunner",
		"uri":    "/api/v3/runners/%s",
	})

	if entity != nil {
		const (
			uri = "/api/v3/runners/%s"
		)
		name := entity.Name.ValueString()

		// Add runner name to span
		if span != nil {
			span.SetData("runner.name", name)
		}

		reqUri := c.uri(format(uri, name))

		// Add API call details to span
		if span != nil {
			span.SetData("api.uri", reqUri)
			span.SetData("api.method", "DELETE")
		}

		kubiyasentry.AddBreadcrumb("client", "Making DELETE request", sentry.LevelInfo, map[string]interface{}{
			"uri":    reqUri,
			"method": "DELETE",
		})

		_, err := c.delete(ctx, reqUri)
		if err != nil {
			kubiyasentry.RecordError(ctx, err)
			kubiyasentry.AddBreadcrumb("client", "Delete request failed", sentry.LevelError, map[string]interface{}{
				"uri":   reqUri,
				"error": err.Error(),
			})
			return err
		}

		// Success breadcrumb
		kubiyasentry.AddBreadcrumb("client", "Successfully deleted runner", sentry.LevelInfo, map[string]interface{}{
			"runner_name": name,
		})

		return nil
	}

	err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
	kubiyasentry.RecordError(ctx, err)
	return err
}

func (c *Client) CreateRunner(ctx context.Context, entity *entities.RunnerModel) (*entities.RunnerModel, error) {
	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateRunner")
		span.SetData("client.resource_type", "runner")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating runner via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateRunner",
		"uri":    "/api/v3/runners/%s",
	})

	if entity == nil {
		err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	const uri = "/api/v3/runners/%s"
	name := entity.Name.ValueString()

	// Add runner name to span
	if span != nil {
		span.SetData("runner.name", name)
	}

	data := struct {
		ManagedBy string `json:"managed_by,omitempty"`
		TaskId    string `json:"task_id,omitempty"`
	}{}
	data.ManagedBy, data.TaskId = managedBy()

	body, err := toJson(data)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "Failed to marshal JSON", sentry.LevelError, map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	reqUri := c.uri(format(uri, name))

	// Add API call details to span
	if span != nil {
		span.SetData("api.uri", reqUri)
		span.SetData("api.method", "POST")
	}

	kubiyasentry.AddBreadcrumb("client", "Making API request", sentry.LevelInfo, map[string]interface{}{
		"uri":    reqUri,
		"method": "POST",
	})

	_, err = c.create(ctx, reqUri, body)
	if err != nil {
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "API request failed", sentry.LevelError, map[string]interface{}{
			"uri":   reqUri,
			"error": err.Error(),
		})
		return nil, err
	}

	// Now call describe to get the runner type
	kubiyasentry.AddBreadcrumb("client", "Reading runner details after creation", sentry.LevelInfo, map[string]interface{}{
		"runner_name": name,
	})

	if err := c.ReadRunner(ctx, entity); err != nil {
		readErr := fmt.Errorf("runner created but failed to read details: %v", err)
		kubiyasentry.RecordError(ctx, readErr)
		kubiyasentry.AddBreadcrumb("client", "Failed to read runner details", sentry.LevelError, map[string]interface{}{
			"runner_name": name,
			"error":       err.Error(),
		})
		return nil, readErr
	}

	// Success - add final breadcrumb
	kubiyasentry.AddBreadcrumb("client", "Successfully created runner", sentry.LevelInfo, map[string]interface{}{
		"runner_name": name,
		"runner_type": entity.RunnerType.ValueString(),
	})

	return entity, nil
}

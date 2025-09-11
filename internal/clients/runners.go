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
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadRunner")
		span.SetData("client.resource_type", "runner")
	}

	if entity == nil {
		logger.Error("ReadRunner called with nil entity")
		err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	runnerName := entity.Name.ValueString()

	logger.Debug("Reading runner",
		"runner_name", runnerName)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading runner via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "ReadRunner",
		"uri":         "/api/v3/runners/%s/describe",
		"runner_name": runnerName,
	})

	const uri = "/api/v3/runners/%s/describe"

	// Add runner name to span
	if span != nil {
		span.SetData("runner.name", runnerName)
	}

	reqUri := c.uri(format(uri, runnerName))

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
		logger.Error("Failed to read runner via API",
			"runner_name", runnerName,
			"uri", reqUri,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "API request failed", sentry.LevelError, map[string]interface{}{
			"uri":   reqUri,
			"error": err.Error(),
		})
		return err
	}

	var r runner
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		logger.Error("Failed to decode runner response",
			"runner_name", runnerName,
			"error", err.Error())
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

	logger.Debug("Runner read successfully",
		"runner_name", r.Name,
		"runner_type", r.RunnerType)

	return nil
}

func (c *Client) DeleteRunner(ctx context.Context, entity *entities.RunnerModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteRunner")
		span.SetData("client.resource_type", "runner")
	}

	if entity == nil {
		logger.Error("DeleteRunner called with nil entity")
		err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	runnerName := entity.Name.ValueString()

	logger.Info("Starting runner deletion",
		"runner_name", runnerName)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting runner via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "DeleteRunner",
		"uri":         "/api/v3/runners/%s",
		"runner_name": runnerName,
	})

	const uri = "/api/v3/runners/%s"

	// Add runner name to span
	if span != nil {
		span.SetData("runner.name", runnerName)
	}

	reqUri := c.uri(format(uri, runnerName))

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
		logger.Error("Failed to delete runner via API",
			"runner_name", runnerName,
			"uri", reqUri,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "Delete request failed", sentry.LevelError, map[string]interface{}{
			"uri":   reqUri,
			"error": err.Error(),
		})
		return err
	}

	// Success breadcrumb
	kubiyasentry.AddBreadcrumb("client", "Successfully deleted runner", sentry.LevelInfo, map[string]interface{}{
		"runner_name": runnerName,
	})

	logger.Info("Runner deleted successfully",
		"runner_name", runnerName)

	return nil
}

func (c *Client) CreateRunner(ctx context.Context, entity *entities.RunnerModel) (*entities.RunnerModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateRunner")
		span.SetData("client.resource_type", "runner")
	}

	if entity == nil {
		logger.Error("CreateRunner called with nil entity")
		err := fmt.Errorf("param entity (*entities.RunnerModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	runnerName := entity.Name.ValueString()

	logger.Info("Starting runner creation",
		"runner_name", runnerName)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating runner via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "CreateRunner",
		"uri":         "/api/v3/runners/%s",
		"runner_name": runnerName,
	})

	const uri = "/api/v3/runners/%s"

	// Add runner name to span
	if span != nil {
		span.SetData("runner.name", runnerName)
	}

	data := struct {
		ManagedBy string `json:"managed_by,omitempty"`
		TaskId    string `json:"task_id,omitempty"`
	}{}
	data.ManagedBy, data.TaskId = managedBy()

	body, err := toJson(data)
	if err != nil {
		logger.Error("Failed to marshal runner data",
			"runner_name", runnerName,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "Failed to marshal JSON", sentry.LevelError, map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	reqUri := c.uri(format(uri, runnerName))

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
		logger.Error("Failed to create runner via API",
			"runner_name", runnerName,
			"uri", reqUri,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.AddBreadcrumb("client", "API request failed", sentry.LevelError, map[string]interface{}{
			"uri":   reqUri,
			"error": err.Error(),
		})
		return nil, err
	}

	// Now call describe to get the runner type
	kubiyasentry.AddBreadcrumb("client", "Reading runner details after creation", sentry.LevelInfo, map[string]interface{}{
		"runner_name": runnerName,
	})

	logger.Info("Runner created, reading details",
		"runner_name", runnerName)

	if err := c.ReadRunner(ctx, entity); err != nil {
		logger.Error("Failed to read runner details after creation",
			"runner_name", runnerName,
			"error", err.Error())
		readErr := fmt.Errorf("runner created but failed to read details: %v", err)
		kubiyasentry.RecordError(ctx, readErr)
		kubiyasentry.AddBreadcrumb("client", "Failed to read runner details", sentry.LevelError, map[string]interface{}{
			"runner_name": runnerName,
			"error":       err.Error(),
		})
		return nil, readErr
	}

	// Success - add final breadcrumb
	kubiyasentry.AddBreadcrumb("client", "Successfully created runner", sentry.LevelInfo, map[string]interface{}{
		"runner_name": runnerName,
		"runner_type": entity.RunnerType.ValueString(),
	})

	logger.Info("Runner created successfully",
		"runner_name", runnerName,
		"runner_type", entity.RunnerType.ValueString())

	return entity, nil
}

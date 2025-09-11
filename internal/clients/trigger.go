package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

// WorkflowDefinition represents the workflow structure for API requests
type WorkflowDefinition struct {
	Name    string                   `json:"name"`
	Version int                      `json:"version"`
	Steps   []map[string]interface{} `json:"steps"`
}

// WorkflowRequest represents the request to create a workflow
type WorkflowRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status"`
	Definition  WorkflowDefinition `json:"definition"`
}

// WorkflowResponse represents the response from workflow creation
type WorkflowResponse struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Definition  WorkflowDefinition `json:"definition"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

// PublishRequest represents the request to publish a workflow with trigger
type PublishRequest struct {
	Type        string `json:"type"`
	WebhookPath string `json:"webhookPath"`
	Runner      string `json:"runner"`
}

// PublishResponse represents the response from publishing a workflow
type PublishResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// WebhookURLRequest represents the request to generate a webhook URL
type WebhookURLRequest struct {
	Runner      string `json:"runner"`
	TriggerType string `json:"triggerType"`
}

// WebhookURLResponse represents the response from webhook URL generation
type WebhookURLResponse struct {
	WebhookUrl  string `json:"webhookUrl"`
	WebhookHash string `json:"webhookHash"`
	Method      string `json:"method"`
}

// normalizeWorkflowJSON normalizes a JSON string to ensure consistent formatting
func normalizeWorkflowJSON(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", nil
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Re-marshal with consistent formatting (compact, no extra spaces)
	normalized, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to normalize JSON: %w", err)
	}

	return string(normalized), nil
}

// parseWorkflowJSON parses the JSON-encoded workflow string to a WorkflowDefinition
func parseWorkflowJSON(workflowJSON string) (*WorkflowDefinition, error) {
	if workflowJSON == "" {
		return nil, fmt.Errorf("workflow JSON is empty")
	}

	// Normalize the JSON first
	normalizedJSON, err := normalizeWorkflowJSON(workflowJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize workflow JSON: %w", err)
	}

	var rawWorkflow map[string]interface{}
	if err := json.Unmarshal([]byte(normalizedJSON), &rawWorkflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
	}

	// Extract name
	name, ok := rawWorkflow["name"].(string)
	if !ok {
		return nil, fmt.Errorf("workflow name is required and must be a string")
	}

	// Extract version
	var version int
	switch v := rawWorkflow["version"].(type) {
	case float64:
		version = int(v)
	case int:
		version = v
	default:
		return nil, fmt.Errorf("workflow version is required and must be a number")
	}

	// Extract steps
	stepsRaw, ok := rawWorkflow["steps"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("workflow steps are required and must be an array")
	}

	steps := make([]map[string]interface{}, 0, len(stepsRaw))
	for _, stepRaw := range stepsRaw {
		step, ok := stepRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("each step must be an object")
		}
		steps = append(steps, step)
	}

	return &WorkflowDefinition{
		Name:    name,
		Version: version,
		Steps:   steps,
	}, nil
}

// workflowDefinitionToJSON converts a WorkflowDefinition to a normalized JSON string
func workflowDefinitionToJSON(def *WorkflowDefinition) (string, error) {
	if def == nil {
		return "", fmt.Errorf("workflow definition is nil")
	}

	jsonBytes, err := json.Marshal(def)
	if err != nil {
		return "", fmt.Errorf("failed to marshal workflow definition: %w", err)
	}

	// Normalize the JSON to ensure consistent formatting
	return normalizeWorkflowJSON(string(jsonBytes))
}

// CreateTrigger creates a new trigger (workflow with webhook)
func (c *Client) CreateTrigger(ctx context.Context, entity *entities.TriggerModel) (*entities.TriggerModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateTrigger")
		span.SetData("client.resource_type", "trigger")
	}

	if entity == nil {
		logger.Error("CreateTrigger called with nil entity")
		err := fmt.Errorf("trigger entity is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	triggerName := entity.Name.ValueString()
	triggerRunner := entity.Runner.ValueString()

	logger.Info("Starting trigger creation",
		"trigger_name", triggerName,
		"runner", triggerRunner)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating trigger via API", sentry.LevelInfo, map[string]interface{}{
		"method":       "CreateTrigger",
		"trigger_name": triggerName,
		"runner":       triggerRunner,
	})

	// Parse workflow JSON to definition
	workflowDef, err := parseWorkflowJSON(entity.Workflow.ValueString())
	if err != nil {
		logger.Error("Failed to parse workflow JSON",
			"trigger_name", triggerName,
			"error", err.Error())
		err = fmt.Errorf("failed to parse workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	// Step 1: Create the workflow in draft status
	workflowReq := WorkflowRequest{
		Name:        entity.Name.ValueString(),
		Description: format("Workflow for trigger %s", entity.Name.ValueString()),
		Status:      "draft",
		Definition:  *workflowDef,
	}

	workflowBody, err := json.Marshal(workflowReq)
	if err != nil {
		logger.Error("Failed to marshal workflow request",
			"trigger_name", triggerName,
			"error", err.Error())
		err = fmt.Errorf("failed to marshal workflow request: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	// Call the workflow creation API
	createPath := "/api/workflows"
	host := "https://composer.kubiya.ai"
	workflowURL := c.uriWithHost(host, createPath)
	resp, err := c.create(ctx, workflowURL, io.NopCloser(strings.NewReader(string(workflowBody))))
	if err != nil {
		logger.Error("Failed to create workflow via API",
			"trigger_name", triggerName,
			"url", workflowURL,
			"error", err.Error())
		err = fmt.Errorf("failed to create workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	var workflowResp WorkflowResponse
	if err := json.NewDecoder(resp).Decode(&workflowResp); err != nil {
		logger.Error("Failed to decode workflow response",
			"trigger_name", triggerName,
			"error", err.Error())
		err = fmt.Errorf("failed to decode workflow response: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Info("Workflow created successfully",
		"trigger_name", triggerName,
		"workflow_id", workflowResp.Id,
		"workflow_status", workflowResp.Status)

	// Step 2: Publish the workflow with webhook trigger
	publishReq := PublishRequest{
		Type:        "webhook",
		WebhookPath: workflowResp.Id,
		Runner:      entity.Runner.ValueString(),
	}

	publishBody, err := json.Marshal(publishReq)
	if err != nil {
		logger.Error("Failed to marshal publish request",
			"trigger_name", triggerName,
			"workflow_id", workflowResp.Id,
			"error", err.Error())
		err = fmt.Errorf("failed to marshal publish request: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	publishPath := format("/api/workflows/%s/publish", workflowResp.Id)
	publishURL := c.uriWithHost(host, publishPath)
	resp, err = c.create(ctx, publishURL, io.NopCloser(strings.NewReader(string(publishBody))))
	if err != nil {
		logger.Error("Failed to publish workflow",
			"trigger_name", triggerName,
			"workflow_id", workflowResp.Id,
			"error", err.Error())
		// Try to clean up the created workflow
		_ = c.deleteWorkflow(ctx, workflowResp.Id)
		err = fmt.Errorf("failed to publish workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Info("Workflow published successfully",
		"trigger_name", triggerName,
		"workflow_id", workflowResp.Id)

	// Step 3: Generate webhook URL
	webhookReq := WebhookURLRequest{
		Runner:      entity.Runner.ValueString(),
		TriggerType: "webhook",
	}

	webhookBody, err := json.Marshal(webhookReq)
	if err != nil {
		logger.Error("Failed to marshal webhook URL request",
			"trigger_name", triggerName,
			"workflow_id", workflowResp.Id,
			"error", err.Error())
		// Try to clean up
		_ = c.deleteWorkflow(ctx, workflowResp.Id)
		err = fmt.Errorf("failed to marshal webhook URL request: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	createTriggerPath := format("/api/workflows/%s/webhook-url", workflowResp.Id)
	webhookURL := c.uriWithHost(host, createTriggerPath)
	resp, err = c.create(ctx, webhookURL, io.NopCloser(strings.NewReader(string(webhookBody))))
	if err != nil {
		logger.Error("Failed to generate webhook URL",
			"trigger_name", triggerName,
			"workflow_id", workflowResp.Id,
			"error", err.Error())
		// Try to clean up
		_ = c.deleteWorkflow(ctx, workflowResp.Id)
		err = fmt.Errorf("failed to generate webhook URL: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	var webhookResp WebhookURLResponse
	if err := json.NewDecoder(resp).Decode(&webhookResp); err != nil {
		logger.Error("Failed to decode webhook URL response",
			"trigger_name", triggerName,
			"workflow_id", workflowResp.Id,
			"error", err.Error())
		// Try to clean up
		_ = c.deleteWorkflow(ctx, workflowResp.Id)
		err = fmt.Errorf("failed to decode webhook URL response: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	// Set the computed fields
	entity.Id = types.StringValue(workflowResp.Id)
	entity.WorkflowId = types.StringValue(workflowResp.Id)
	entity.Url = types.StringValue(webhookResp.WebhookUrl)
	entity.Status = types.StringValue("published")

	// Normalize the workflow JSON to ensure consistency
	normalizedWorkflow, err := workflowDefinitionToJSON(&workflowResp.Definition)
	if err != nil {
		// If we can't normalize, keep the original
		normalizedWorkflow = entity.Workflow.ValueString()
	}
	entity.Workflow = types.StringValue(normalizedWorkflow)

	logger.Info("Trigger created successfully",
		"trigger_name", triggerName,
		"trigger_id", workflowResp.Id,
		"webhook_url", webhookResp.WebhookUrl)

	return entity, nil
}

// ReadTrigger reads an existing trigger
func (c *Client) ReadTrigger(ctx context.Context, entity *entities.TriggerModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadTrigger")
		span.SetData("client.resource_type", "trigger")
	}

	if entity == nil {
		logger.Error("ReadTrigger called with nil entity")
		err := fmt.Errorf("trigger entity is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	workflowId := entity.WorkflowId.ValueString()
	if workflowId == "" {
		workflowId = entity.Id.ValueString()
	}
	triggerName := entity.Name.ValueString()

	logger.Debug("Reading trigger",
		"trigger_name", triggerName,
		"workflow_id", workflowId)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading trigger via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "ReadTrigger",
		"workflow_id": workflowId,
	})

	// Get workflow details
	host := "https://composer.kubiya.ai"
	readPath := format("/api/workflows/%s", workflowId)
	workflowURL := c.uriWithHost(host, readPath)

	resp, err := c.read(ctx, workflowURL)
	if err != nil {
		logger.Error("Failed to read workflow",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to read workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	var workflowResp WorkflowResponse
	if err := json.NewDecoder(resp).Decode(&workflowResp); err != nil {
		logger.Error("Failed to decode workflow response",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to decode workflow response: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	// Update entity with current state
	entity.Status = types.StringValue(workflowResp.Status)
	entity.WorkflowId = types.StringValue(workflowResp.Id)

	// Store the normalized workflow definition to ensure consistency
	normalizedWorkflow, err := workflowDefinitionToJSON(&workflowResp.Definition)
	if err != nil {
		// If we can't normalize, keep the existing value
		normalizedWorkflow = entity.Workflow.ValueString()
	}
	entity.Workflow = types.StringValue(normalizedWorkflow)

	logger.Debug("Trigger read successfully",
		"trigger_name", triggerName,
		"workflow_id", workflowResp.Id,
		"status", workflowResp.Status)

	// The webhook URL should remain the same, so we don't update it here

	return nil
}

// UpdateTrigger updates an existing trigger
func (c *Client) UpdateTrigger(ctx context.Context, entity *entities.TriggerModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "UpdateTrigger")
		span.SetData("client.resource_type", "trigger")
	}

	if entity == nil {
		logger.Error("UpdateTrigger called with nil entity")
		err := fmt.Errorf("trigger entity is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	workflowId := entity.WorkflowId.ValueString()
	if workflowId == "" {
		workflowId = entity.Id.ValueString()
	}
	triggerName := entity.Name.ValueString()

	logger.Info("Starting trigger update",
		"trigger_name", triggerName,
		"workflow_id", workflowId)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Updating trigger via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "UpdateTrigger",
		"workflow_id": workflowId,
	})

	// Parse workflow JSON to definition
	workflowDef, err := parseWorkflowJSON(entity.Workflow.ValueString())
	if err != nil {
		logger.Error("Failed to parse workflow JSON during update",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to parse workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	// Update the workflow
	workflowReq := WorkflowRequest{
		Name:        entity.Name.ValueString(),
		Description: format("Workflow for trigger %s", entity.Name.ValueString()),
		Status:      "published",
		Definition:  *workflowDef,
	}

	workflowBody, err := json.Marshal(workflowReq)
	if err != nil {
		logger.Error("Failed to marshal workflow request during update",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to marshal workflow request: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	host := "https://composer.kubiya.ai"
	updatePath := format("/api/workflows/%s", workflowId)
	workflowURL := c.uriWithHost(host, updatePath)

	resp, err := c.update(ctx, workflowURL, io.NopCloser(strings.NewReader(string(workflowBody))))
	if err != nil {
		logger.Error("Failed to update workflow via API",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to update workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	var workflowResp WorkflowResponse
	if err := json.NewDecoder(resp).Decode(&workflowResp); err != nil {
		logger.Error("Failed to decode workflow response during update",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		err = fmt.Errorf("failed to decode workflow response: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	// If runner changed, we might need to regenerate the webhook URL
	// For now, we'll keep the existing URL

	entity.Status = types.StringValue(workflowResp.Status)

	logger.Info("Trigger updated successfully",
		"trigger_name", triggerName,
		"workflow_id", workflowId,
		"status", workflowResp.Status)

	return nil
}

// DeleteTrigger deletes an existing trigger
func (c *Client) DeleteTrigger(ctx context.Context, entity *entities.TriggerModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteTrigger")
		span.SetData("client.resource_type", "trigger")
	}

	if entity == nil {
		logger.Error("DeleteTrigger called with nil entity")
		err := fmt.Errorf("trigger entity is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	workflowId := entity.WorkflowId.ValueString()
	if workflowId == "" {
		workflowId = entity.Id.ValueString()
	}
	triggerName := entity.Name.ValueString()

	logger.Info("Starting trigger deletion",
		"trigger_name", triggerName,
		"workflow_id", workflowId)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting trigger via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "DeleteTrigger",
		"workflow_id": workflowId,
	})

	err := c.deleteWorkflow(ctx, workflowId)
	if err != nil {
		logger.Error("Failed to delete trigger",
			"trigger_name", triggerName,
			"workflow_id", workflowId,
			"error", err.Error())
		return err
	}

	logger.Info("Trigger deleted successfully",
		"trigger_name", triggerName,
		"workflow_id", workflowId)

	return nil
}

// deleteWorkflow is a helper function to delete a workflow
func (c *Client) deleteWorkflow(ctx context.Context, workflowId string) error {
	host := "https://composer.kubiya.ai"
	deletePath := format("/api/workflows/%s", workflowId)
	workflowURL := c.uriWithHost(host, deletePath)
	resp, err := c.delete(ctx, workflowURL)
	if err != nil {
		err = fmt.Errorf("failed to delete workflow: %w", err)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	// Check if deletion was successful
	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp).Decode(&result); err != nil {
		// If we can't decode, assume success if no error was returned
		return nil
	}

	if !result.Success {
		err := fmt.Errorf("failed to delete workflow: %s", result.Message)
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	return nil
}

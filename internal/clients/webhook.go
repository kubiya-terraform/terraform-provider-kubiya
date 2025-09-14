package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type (
	webhook struct {
		Id            string         `json:"id"`
		Name          string         `json:"name"`
		Filter        string         `json:"filter"`
		Prompt        string         `json:"prompt"`
		Source        string         `json:"source"`
		AgentId       string         `json:"agent_id"`
		CreatedAt     time.Time      `json:"created_at"`
		CreatedBy     string         `json:"created_by"`
		UpdatedAt     time.Time      `json:"updated_at"`
		WebhookUrl    string         `json:"webhook_url"`
		TaskId        string         `json:"task_id"`
		ManagedBy     string         `json:"managed_by"`
		Communication *communication `json:"communication"`

		Runner   string `json:"runner,omitempty"`
		Workflow string `json:"workflow,omitempty"`
	}

	communication struct {
		Method      string `json:"method"`
		Destination string `json:"destination"` // prefix # = channel, @ = person (lookup for his email)
	}
)

func toWebhook(w *entities.WebhookModel, cs *state) (*webhook, error) {
	wh := &webhook{
		Id:         w.Id.ValueString(),
		WebhookUrl: w.Url.ValueString(),
		Name:       w.Name.ValueString(),
		Filter:     w.Filter.ValueString(),
		Prompt:     w.Prompt.ValueString(),
		Source:     w.Source.ValueString(),
		CreatedBy:  w.CreatedBy.ValueString(),
		Runner:     w.Runner.ValueString(),
		Workflow:   w.Workflow.ValueString(),
	}

	var tmp interface{}
	if wh.Workflow != "" && wh.Workflow != "{}" {
		if err := json.Unmarshal([]byte(wh.Workflow), &tmp); err != nil {
			return nil, err
		}
	}

	data := w.Workflow.ValueString()
	normalized, err := normalizeJSON(string(data))
	if err != nil {
		return wh, err
	}

	if normalized == "[]" {
		normalized = ""
	}
	wh.Workflow = normalized

	for _, a := range cs.agentList {
		if equal(a.Name, w.Agent.ValueString()) {
			wh.AgentId = a.Uuid
			break
		}
	}

	// Get method, default to "Slack" with capital S if empty
	method := w.Method.ValueString()
	if method == "" {
		method = "Slack" // Capital S for consistency
	}

	// Handle destination based on method
	if strings.EqualFold(method, "http") {
		// For http, destination can be empty
		wh.Communication = &communication{Method: method, Destination: ""}
	} else if len(w.Destination.ValueString()) >= 1 {
		const (
			at    = "@"
			pound = "#"
		)

		destination := w.Destination.ValueString()

		if !strings.HasPrefix(destination, pound) {
			t := strings.TrimPrefix(destination, at)
			for _, u := range cs.userList {
				if equal(t, u.Name) {
					destination = u.Email
					break
				}
			}
		}

		// Special handling for teams method
		if strings.EqualFold(method, "teams") {
			teamName := w.TeamName.ValueString()
			channelName := strings.TrimPrefix(destination, pound)
			destination = fmt.Sprintf("#{\"team_name\":\"%s\",\"channel_name\":\"%s\"}",
				teamName, channelName)
		}

		wh.Communication = &communication{Method: method, Destination: destination}
	}

	if wh.Communication == nil {
		wh.Communication = &communication{Method: "http", Destination: ""}
	} else {
		if wh.Communication.Method == "" {
			wh.Communication.Method = "http"
		}
		if wh.Communication.Destination == "" {
			wh.Communication.Destination = "webhook"
		}
	}

	return wh, nil
}

func fromWebhook(w *webhook, cs *state) (*entities.WebhookModel, error) {
	agentName := ""
	destination := ""
	at := w.CreatedAt.String()

	if w.Communication != nil {
		destination = w.Communication.Destination
	}

	for _, a := range cs.agentList {
		if strings.EqualFold(w.AgentId, a.Uuid) {
			agentName = a.Name
			break
		}
	}

	wf, err := normalizeJSON(w.Workflow)
	if err != nil {
		return nil, err
	}

	if wf == "[]" {
		wf = ""
	}

	wh := &entities.WebhookModel{
		CreatedAt:   types.StringValue(at),
		Id:          types.StringValue(w.Id),
		Name:        types.StringValue(w.Name),
		Filter:      types.StringValue(w.Filter),
		Source:      types.StringValue(w.Source),
		Prompt:      types.StringValue(w.Prompt),
		Agent:       types.StringValue(agentName),
		CreatedBy:   types.StringValue(w.CreatedBy),
		Destination: types.StringValue(destination),
		Url:         types.StringValue(w.WebhookUrl),
		Workflow:    types.StringValue(wf),
		Runner:      types.StringValue(w.Runner),
	}

	// Set method and team_name fields
	if w.Communication != nil {
		wh.Method = types.StringValue(w.Communication.Method)

		// For teams method, extract team_name from the destination JSON
		if strings.EqualFold(w.Communication.Method, "teams") &&
			strings.HasPrefix(w.Communication.Destination, "#{") {
			// Remove the "#" prefix
			jsonStr := strings.TrimPrefix(w.Communication.Destination, "#")
			var teamsDest struct {
				TeamName    string `json:"team_name"`
				ChannelName string `json:"channel_name"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &teamsDest); err == nil {
				wh.TeamName = types.StringValue(teamsDest.TeamName)
				wh.Destination = types.StringValue(teamsDest.ChannelName) // Don't add # prefix
			}
		}
	}

	return wh, nil
}

func (c *Client) ReadWebhook(ctx context.Context, entity *entities.WebhookModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if entity != nil {
		id := entity.Id.ValueString()
		name := entity.Name.ValueString()

		logger.Debug("Reading webhook", map[string]interface{}{
			"webhook_id":   id,
			"webhook_name": name,
		})

		// Add breadcrumb for operation tracking
		kubiyasentry.AddBreadcrumb("client", "Reading webhook", sentry.LevelDebug, map[string]interface{}{
			"method":       "ReadWebhook",
			"webhook_id":   id,
			"webhook_name": name,
		})

		cs, err := c.state()
		if err != nil {
			logger.Error("Failed to get client state", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}

		for _, w := range cs.webhookList {
			if equal(w.Id, id) || equal(w.Name, name) {
				entity, err = fromWebhook(w, cs)
				if err != nil {
					logger.Error("Failed to convert webhook from API response", map[string]interface{}{
						"error":        err.Error(),
						"webhook_id":   id,
						"webhook_name": name,
					})
					return err
				}

				logger.Debug("Successfully read webhook", map[string]interface{}{
					"webhook_id":   id,
					"webhook_name": name,
				})
				break
			}
		}

		return err
	}

	err := fmt.Errorf("param entity (*entities.WebhookModel) is nil")
	logger.Error("Webhook entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	return err
}

func (c *Client) DeleteWebhook(ctx context.Context, entity *entities.WebhookModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if entity != nil {
		const (
			ok     = ""
			path   = "/api/v1/event/%s"
			errMsg = "failed to delete webhook - %s"
		)

		id := entity.Id.ValueString()
		name := entity.Name.ValueString()

		logger.Info("Starting webhook deletion", map[string]interface{}{
			"webhook_id":   id,
			"webhook_name": name,
		})

		// Add breadcrumb for operation tracking
		kubiyasentry.AddBreadcrumb("client", "Deleting webhook", sentry.LevelInfo, map[string]interface{}{
			"method":       "DeleteWebhook",
			"webhook_id":   id,
			"webhook_name": name,
		})

		uri := c.uri(fmt.Sprintf(path, id))
		resp, err := c.delete(ctx, uri)
		if err != nil {
			logger.Error("Failed to delete webhook", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		r := &struct {
			Result string `json:"result"`
		}{}

		err = json.NewDecoder(resp).Decode(&r)
		if err != nil || r == nil {
			if err != nil {
				logger.Error("Failed to decode delete webhook response", map[string]interface{}{
					"error":        err.Error(),
					"webhook_id":   id,
					"webhook_name": name,
				})
				return err
			}
			err = fmt.Errorf(errMsg, id)
			logger.Error("Invalid delete webhook response", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		if strings.Contains(r.Result, ok) {
			logger.Info("Successfully deleted webhook", map[string]interface{}{
				"webhook_id":   id,
				"webhook_name": name,
			})
			return nil
		}

		err = fmt.Errorf(errMsg, id)
		logger.Error("Failed to delete webhook", map[string]interface{}{
			"error":        err.Error(),
			"webhook_id":   id,
			"webhook_name": name,
			"result":       r.Result,
		})
		return err
	}

	err := fmt.Errorf("param entity (*entities.WebhookModel) is nil")
	logger.Error("Webhook entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	return err
}

func (c *Client) UpdateWebhook(ctx context.Context, entity *entities.WebhookModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if entity != nil {
		id := entity.Id.ValueString()
		name := entity.Name.ValueString()
		wf := entity.Workflow.ValueString()
		agentId := entity.Agent.ValueString()

		logger.Info("Starting webhook update", map[string]interface{}{
			"webhook_id":   id,
			"webhook_name": name,
		})

		// Add breadcrumb for operation tracking
		kubiyasentry.AddBreadcrumb("client", "Updating webhook", sentry.LevelInfo, map[string]interface{}{
			"method":       "UpdateWebhook",
			"webhook_id":   id,
			"webhook_name": name,
		})

		if (wf == "" && agentId == "") || (wf != "" && agentId != "") {
			err := fmt.Errorf("workflow or agent is required")
			logger.Error("Invalid webhook configuration", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
				"has_workflow": wf != "",
				"has_agent":    agentId != "",
			})
			return err
		}

		const (
			path = "/api/v1/event/%s"
		)

		cs, err := c.state()
		if err != nil {
			logger.Error("Failed to get client state", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		uri := c.uri(format(path, id))

		data, err := toWebhook(entity, cs)
		if err != nil {
			logger.Error("Failed to convert webhook to API format", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}
		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			logger.Error("Failed to marshal webhook data", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to update webhook", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		var r *webhook
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode update webhook response", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
			return err
		}

		entity, err = fromWebhook(r, cs)
		if err != nil {
			logger.Error("Failed to convert updated webhook from API response", map[string]interface{}{
				"error":        err.Error(),
				"webhook_id":   id,
				"webhook_name": name,
			})
		} else {
			logger.Info("Successfully updated webhook", map[string]interface{}{
				"webhook_id":   id,
				"webhook_name": name,
			})
		}

		return err
	}

	err := fmt.Errorf("param entity (*entities.WebhookModel) is nil")
	logger.Error("Webhook entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	return err
}

func (c *Client) CreateWebhook(ctx context.Context, entity *entities.WebhookModel) (*entities.WebhookModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateWebhook")
		span.SetData("client.resource_type", "webhook")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating webhook via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateWebhook",
		"uri":    "/api/v1/event",
	})

	if entity != nil {
		name := entity.Name.ValueString()
		wf := entity.Workflow.ValueString()
		agentId := entity.Agent.ValueString()

		logger.Info("Starting webhook creation", map[string]interface{}{
			"webhook_name": name,
		})

		if (wf == "" && agentId == "") || (wf != "" && agentId != "") {
			err := fmt.Errorf("workflow or agent is required")
			logger.Error("Invalid webhook configuration", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
				"has_workflow": wf != "",
				"has_agent":    agentId != "",
			})
			return nil, err
		}

		cs, err := c.state()
		if err != nil {
			logger.Error("Failed to get client state", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
			return nil, err
		}

		uri := c.uri("/api/v1/event")

		data, err := toWebhook(entity, cs)
		if err != nil {
			logger.Error("Failed to convert webhook to API format", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
			return nil, err
		}
		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			logger.Error("Failed to marshal webhook data", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
			return nil, err
		}

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to create webhook", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
			return nil, err
		}

		var r *webhook
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode create webhook response", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
			return nil, err
		}

		result, err := fromWebhook(r, cs)
		if err != nil {
			logger.Error("Failed to convert created webhook from API response", map[string]interface{}{
				"error":        err.Error(),
				"webhook_name": name,
			})
		} else {
			logger.Info("Successfully created webhook", map[string]interface{}{
				"webhook_id":   result.Id.ValueString(),
				"webhook_name": name,
			})
		}

		return result, err
	}

	err := fmt.Errorf("param entity (*entities.WebhookModel) is nil")
	logger.Error("Webhook entity is nil", map[string]interface{}{
		"error": err.Error(),
	})
	return nil, err
}

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
	}

	communication struct {
		Method      string `json:"method"`
		Destination string `json:"destination"` // prefix # = channel, @ = person (lookup for his email)
	}
)

func toWebhook(w *entities.WebhookModel, cs *state) *webhook {
	wh := &webhook{
		Id:         w.Id.ValueString(),
		WebhookUrl: w.Url.ValueString(),
		Name:       w.Name.ValueString(),
		Filter:     w.Filter.ValueString(),
		Prompt:     w.Prompt.ValueString(),
		Source:     w.Source.ValueString(),
	}

	for _, a := range cs.agentList {
		if equal(a.Name, w.Agent.ValueString()) {
			wh.AgentId = a.Uuid
			break
		}
	}

	if len(w.Destination.ValueString()) >= 1 {
		const (
			at     = "@"
			pound  = "#"
			Method = "Slack"
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

		wh.Communication = &communication{Method: Method, Destination: destination}
	}

	return wh
}

func fromWebhook(w *webhook, cs *state) *entities.WebhookModel {
	by := ""
	agentName := ""
	destination := ""
	at := w.CreatedAt.String()

	if w.Communication != nil {
		destination = w.Communication.Destination
	}

	for _, u := range cs.userList {
		if strings.EqualFold(w.CreatedBy, u.UUID) {
			by = u.Email
			break
		}
	}

	for _, a := range cs.agentList {
		if strings.EqualFold(w.AgentId, a.Uuid) {
			agentName = a.Name
			break
		}
	}

	wh := &entities.WebhookModel{
		CreatedAt:   types.StringValue(at),
		CreatedBy:   types.StringValue(by),
		Id:          types.StringValue(w.Id),
		Name:        types.StringValue(w.Name),
		Filter:      types.StringValue(w.Filter),
		Source:      types.StringValue(w.Source),
		Prompt:      types.StringValue(w.Prompt),
		Agent:       types.StringValue(agentName),
		Destination: types.StringValue(destination),
	}

	return wh
}

func (c *Client) ReadWebhook(_ context.Context, entity *entities.WebhookModel) error {
	if entity != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := entity.Id.ValueString()
		name := entity.Name.ValueString()

		for _, w := range cs.webhookList {
			if equal(w.Id, id) || equal(w.Name, name) {
				entity = fromWebhook(w, cs)
				break
			}
		}

		return err
	}

	return fmt.Errorf("param entity (*entities.WebhookModel) is nil")
}

func (c *Client) DeleteWebhook(ctx context.Context, entity *entities.WebhookModel) error {
	if entity != nil {
		const (
			ok     = ""
			path   = "/api/v1/event/%s"
			errMsg = "failed to delete webhook - %s"
		)

		id := entity.Id.ValueString()
		uri := c.uri(fmt.Sprintf(path, id))
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
			return fmt.Errorf(errMsg, id)
		}

		if strings.Contains(r.Result, ok) {
			return nil
		}

		return fmt.Errorf(errMsg, id)
	}

	return fmt.Errorf("param entity (*entities.WebhookModel) is nil")
}

func (c *Client) UpdateWebhook(ctx context.Context, entity *entities.WebhookModel) error {
	if entity != nil {
		const (
			path = "/api/v1/event/%s"
		)

		cs, err := c.state()
		if err != nil {
			return err
		}

		id := entity.Id.ValueString()

		uri := c.uri(format(path, id))

		data := toWebhook(entity, cs)
		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			return err
		}

		var r *webhook
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return err
		}

		entity = fromWebhook(r, cs)

		return err
	}
	return fmt.Errorf("param entity (*entities.WebhookModel) is nil")
}

func (c *Client) CreateWebhook(ctx context.Context, entity *entities.WebhookModel) (*entities.WebhookModel, error) {
	if entity != nil {
		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/event")

		data := toWebhook(entity, cs)
		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		var r *webhook
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		return fromWebhook(r, cs), err
	}

	return nil, fmt.Errorf("param entity (*entities.WebhookModel) is nil")
}

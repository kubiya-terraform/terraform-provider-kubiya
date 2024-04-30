package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type (
	webhook struct {
		Id            string         `json:"id"`
		Org           string         `json:"org"`
		Name          string         `json:"name"`
		Filter        string         `json:"filter"`
		Prompt        string         `json:"prompt"`
		Source        string         `json:"source"`
		AgentId       string         `json:"agent_id"`
		CreatedAt     time.Time      `json:"created_at"`
		CreatedBy     string         `json:"created_by"`
		UpdatedAt     time.Time      `json:"updated_at"`
		WebhookUrl    string         `json:"webhook_url"`
		Communication *communication `json:"communication"`
	}

	communication struct {
		Method      string `json:"method"`
		Destination string `json:"destination"` // prefix # = channel, @ = person (lookup for his email)
	}
)

func (c *Client) DeleteWebhook(req *entities.WebhookModel) error {
	const (
		method = "DELETE"
		uri    = "/api/v1/event/%s"
		errMsg = "failed to delete webhook - %s"
	)

	id := req.Id.ValueString()

	reqUri := fmt.Sprintf(uri, id)

	respBody, err := c.webhook(method, reqUri, nil)
	if err != nil || respBody == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf(errMsg, id)
	}

	result := &struct {
		Result string `json:"result"`
	}{}

	err = json.NewDecoder(respBody).Decode(&result)
	if err != nil || result == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf(errMsg, id)
	}

	if strings.Contains(result.Result, "ok") {
		return nil
	}

	return fmt.Errorf(errMsg, id)
}

func (c *Client) webhook(m, uri string, body io.Reader) (io.Reader, error) {
	const (
		slash  = "/"
		layout = "%s/%s"
	)

	for {
		if !strings.HasPrefix(uri, slash) {
			break
		}

		uri = strings.TrimPrefix(uri, slash)
	}

	uri = fmt.Sprintf(layout, c.host, uri)

	req, err := http.NewRequest(m, uri, body)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	return c.doReaderHttpRequest(req)

}

func (c *Client) GetWebhook(req *entities.WebhookModel) (*entities.WebhookModel, error) {
	const (
		method = "GET"
		uri    = "/api/v1/event"
		errMsg = "failed to get webhook - %s"
	)

	id := req.Id.ValueString()
	name := req.Name.ValueString()

	respBody, err := c.webhook(method, uri, nil)
	if err != nil || respBody == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errMsg, name)
	}

	var result []*webhook

	if err = json.NewDecoder(respBody).Decode(&result); err != nil {
		return nil, err
	}

	if len(result) <= 0 {
		return nil, fmt.Errorf(errMsg, name)
	}

	users, err := c.users()
	if err != nil {
		return nil, err
	}

	agents, err := c.agents()
	if err != nil {
		return nil, err
	}

	var wbResp *entities.WebhookModel

	for _, i := range result {
		if strings.EqualFold(i.Id, id) ||
			strings.EqualFold(i.Name, name) {
			by := ""
			agentName := ""
			destination := ""
			at := i.CreatedAt.String()

			for _, u := range users {
				if strings.EqualFold(i.CreatedBy, u.UUID) {
					by = u.Email
					break
				}
			}

			for _, a := range agents {
				if strings.EqualFold(i.AgentId, a.Uuid) {
					agentName = a.Name
					break
				}
			}

			if i.Communication != nil {
				destination = i.Communication.Destination
			}

			wbResp = &entities.WebhookModel{
				CreatedAt:   types.StringValue(at),
				CreatedBy:   types.StringValue(by),
				Id:          types.StringValue(i.Id),
				Name:        types.StringValue(i.Name),
				Filter:      types.StringValue(i.Filter),
				Source:      types.StringValue(i.Source),
				Prompt:      types.StringValue(i.Prompt),
				Agent:       types.StringValue(agentName),
				Destination: types.StringValue(destination),
			}

			break
		}
	}

	return wbResp, err
}

func (c *Client) CreateWebhook(req *entities.WebhookModel) (*entities.WebhookModel, error) {
	const (
		method = "POST"
		uri    = "/api/v1/event"
		errMsg = "failed to create webhook"
	)

	wbReq := &webhook{
		Name:   req.Name.ValueString(),
		Filter: req.Filter.ValueString(),
		Prompt: req.Prompt.ValueString(),
		Source: req.Source.ValueString(),
	}

	me, err := c.self()
	if err != nil {
		return nil, err
	}
	wbReq.CreatedBy = me.Email

	if len(req.Destination.ValueString()) >= 1 {
		const (
			at     = "@"
			pound  = "#"
			Method = "Slack"
		)

		users, e := c.users()
		if e != nil {
			return nil, e
		}

		destination := req.Destination.ValueString()

		if !strings.HasPrefix(destination, pound) {
			t := strings.TrimPrefix(destination, at)
			for _, userItem := range users {
				if strings.EqualFold(t, userItem.Name) {
					destination = userItem.Email
					break
				}
			}
		}

		wbReq.Communication = &communication{Method: Method, Destination: destination}
	}

	agents, err := c.agents()
	if err != nil {
		return nil, err
	}

	for _, agent := range agents {
		if agent.Name == req.Agent.ValueString() {
			wbReq.AgentId = agent.Uuid
		}
	}

	reqBody, err := toJson(wbReq)
	if err != nil || reqBody == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("webhook is nil")
	}

	respBody, err := c.webhook(method, uri, reqBody)
	if err != nil || respBody == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errMsg)
	}

	var wbResp webhook

	if err = json.NewDecoder(respBody).Decode(&wbResp); err != nil {
		return nil, err
	}

	var agentName string
	var destination string

	for _, agent := range agents {
		if agent.Name == req.Agent.ValueString() {
			agentName = agent.Name
		}
	}

	if wbResp.Communication != nil {
		destination = wbResp.Communication.Destination
	}

	return &entities.WebhookModel{
		Agent:       types.StringValue(agentName),
		Id:          types.StringValue(wbResp.Id),
		Name:        types.StringValue(wbResp.Name),
		Destination: types.StringValue(destination),
		Filter:      types.StringValue(wbResp.Filter),
		Source:      types.StringValue(wbResp.Source),
		Prompt:      types.StringValue(wbResp.Prompt),
		CreatedBy:   types.StringValue(wbResp.CreatedBy),
		CreatedAt:   types.StringValue(wbResp.CreatedAt.String()),
	}, nil
}

func (c *Client) UpdateWebhook(req *entities.WebhookModel) (*entities.WebhookModel, error) {
	const (
		method = "PUT"
		uri    = "/api/v1/event/%s"
		errMsg = "failed to update webhook - %s"
	)

	wbReq := &webhook{
		Name:   req.Name.ValueString(),
		Filter: req.Filter.ValueString(),
		Prompt: req.Prompt.ValueString(),
		Source: req.Source.ValueString(),
	}

	id := req.Id.ValueString()

	me, err := c.self()
	if err != nil {
		return nil, err
	}
	wbReq.CreatedBy = me.Email

	if len(req.Destination.ValueString()) >= 1 {
		const (
			at     = "@"
			pound  = "#"
			Method = "Slack"
		)

		users, e := c.users()
		if e != nil {
			return nil, e
		}

		destination := req.Destination.ValueString()

		if !strings.HasPrefix(destination, pound) {
			t := strings.TrimPrefix(destination, at)
			for _, userItem := range users {
				if strings.EqualFold(t, userItem.Name) {
					destination = userItem.Email
					break
				}
			}
		}

		wbReq.Communication = &communication{Method: Method, Destination: destination}
	}

	agents, err := c.agents()
	if err != nil {
		return nil, err
	}

	for _, agent := range agents {
		if agent.Name == req.Agent.ValueString() {
			wbReq.AgentId = agent.Uuid
		}
	}

	reqBody, err := toJson(wbReq)
	if err != nil || reqBody == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errMsg, id)
	}

	reqUri := fmt.Sprintf(uri, id)

	respBody, err := c.webhook(method, reqUri, reqBody)
	if err != nil || respBody == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errMsg, id)
	}

	var wbResp webhook

	if err = json.NewDecoder(respBody).Decode(&wbResp); err != nil {
		return nil, err
	}

	var agentName string
	var destination string

	for _, agent := range agents {
		if agent.Name == req.Agent.ValueString() {
			agentName = agent.Name
		}
	}

	if wbResp.Communication != nil {
		destination = wbResp.Communication.Destination
	}

	return &entities.WebhookModel{
		Agent:       types.StringValue(agentName),
		Id:          types.StringValue(wbResp.Id),
		Name:        types.StringValue(wbResp.Name),
		Destination: types.StringValue(destination),
		Filter:      types.StringValue(wbResp.Filter),
		Source:      types.StringValue(wbResp.Source),
		Prompt:      types.StringValue(wbResp.Prompt),
		CreatedBy:   types.StringValue(wbResp.CreatedBy),
		CreatedAt:   types.StringValue(wbResp.CreatedAt.String()),
	}, nil
}

package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Client struct {
	host    string
	userKey string
	client  *http.Client
}

func New(uk string) (*Client, error) {
	if len(uk) >= 1 {
		client := &http.Client{}
		host := "https://api.kubiya.ai"
		return &Client{userKey: uk, client: client, host: host}, nil
	}

	return nil, eformat("ApiKey is missing or empty")
}

func (c *Client) self() (*user, error) {
	const (
		path = "/api/v1/users/self"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result *user
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) state() (*state, error) {
	var err error
	var currentState state

	if users, e := c.users(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.users = append(make([]*user, 0), users...)
	}

	if agents, e := c.agents(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.agents = append(make([]*agent, 0), agents...)
	}

	if groups, e := c.groups(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.groups = append(make([]*group, 0), groups...)
	}

	if models, e := c.models(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.models = append(make([]string, 0), models...)
	}

	if runners, e := c.runners(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.runners = append(make([]*runner, 0), runners...)
	}

	if secrets, e := c.secrets(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.secrets = append(make([]*secret, 0), secrets...)
	}

	if webhooks, e := c.webhooks(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.webhooks = append(make([]*webhook, 0), webhooks...)
	}

	if integrations, e := c.integrations(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.integrations = append(make([]*integration, 0), integrations...)
	}

	return &currentState, err
}

func (c *Client) users() ([]*user, error) {
	const (
		path = "/api/v1/users"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*user
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) agents() ([]*agent, error) {
	const (
		path = "/api/v1/agents"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*agent
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) groups() ([]*group, error) {
	const (
		path = "/api/v1/manage/groups"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.readBytes(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*group
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&result)

	return result, err
}

func (c *Client) models() ([]string, error) {
	const (
		sep  = ","
		path = "/api/v1/featureflags"
		body = `["supported_llm_models"]`
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.create(ctx, uri,
		strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	tmp := &struct {
		Models string `json:"supported_llm_models"`
	}{}

	if err = json.NewDecoder(resp).Decode(tmp); err != nil {
		return nil, err
	}

	var result []string

	for _, item := range strings.Split(tmp.Models, sep) {
		result = append(result, strings.TrimSpace(item))
	}

	return result, err
}

func (c *Client) runners() ([]*runner, error) {
	const (
		path = "/api/v3/runners"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*runner
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) secrets() ([]*secret, error) {
	const (
		path = "/api/v1/secrets"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*secret
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) webhooks() ([]*webhook, error) {
	const (
		path = "/api/v1/event"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var result []*webhook
	err = json.NewDecoder(resp).Decode(&result)

	return result, err
}

func (c *Client) integrations() ([]*integration, error) {
	const (
		path            = "/api/v1/runners"
		managed         = "kubiya-managed"
		pathIntegration = "/api/v2/integrations"
	)

	uri := c.uri(path)
	ctx := context.Background()

	resp, err := c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	var tmp map[string]interface{}

	err = json.NewDecoder(resp).Decode(&tmp)
	if err != nil {
		return nil, err
	}

	result := []*integration{
		{Name: "slack"},
	}

	if val, ok := tmp[managed]; ok {
		if items, ok := val.(map[string]interface{}); ok {
			for integrationName := range items {
				result = append(result, &integration{
					Name: integrationName,
				})
			}
		}
	}

	resp, err = c.read(ctx,
		c.uri(pathIntegration))
	if err != nil {
		return nil, err
	}

	var tmpList []*integrationApi
	err = json.NewDecoder(resp).Decode(&tmpList)
	if err != nil {
		return nil, err
	}

	for _, item := range tmpList {
		result = append(result, &integration{
			Name: item.Name,
		})
	}

	return result, err
}

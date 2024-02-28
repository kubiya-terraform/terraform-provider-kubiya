package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultHost    = "https://api.kubiya.ai"
	userKeyError   = "UserKey is empty or nil"
	defaultTimeout = 10 * time.Second
)

type Client struct {
	host         string
	email        string
	userKey      string
	organization string
	client       *http.Client
}

func NewClient(key, email, organization string) (*Client, error) {
	if len(key) >= 1 {
		host := defaultHost
		timeout := defaultTimeout
		client := &http.Client{Timeout: timeout}
		return &Client{host: host, email: email,
			userKey: key, organization: organization, client: client}, nil
	}

	return nil, fmt.Errorf(userKeyError)
}

func (c *Client) queryParams(uri string) string {
	if len(c.organization) >= 1 && len(c.email) >= 1 {
		t := "%s?organization=%s&email=%s"
		return fmt.Sprintf(t, uri, c.organization, c.email)
	}

	return uri
}

// DeleteAgent DELETE /api/v1/agents/{id}
// https://api.kubiya.ai/api/v1/agents/01b81e08-17eb-4a3e-b0c6-6a48b0f2fad0
func (c *Client) DeleteAgent(id string) error {
	m := "DELETE"
	t := "%s/api/v1/agents/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, id))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to create *http.Request")
	}

	if _, err = c.doBytesHttpRequest(req); err != nil {
		return err
	}

	return nil
}

// GetAgents GET /api/v1/agents
func (c *Client) GetAgents() ([]*Agent, error) {
	m := "GET"
	t := "%s/api/v1/agents"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result []*Agent
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteRunner DELETE /api/v1/runners/{name}
func (c *Client) DeleteRunner(name string) error {
	m := "DELETE"
	t := "%s/api/v1/runners/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, name))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to create *http.Request")
	}

	if _, err = c.doBytesHttpRequest(req); err != nil {
		return err
	}

	return nil
}

// GetAgentById GET /api/v1/agents/{id}
func (c *Client) GetAgentById(id string) (*Agent, error) {
	m := "GET"
	t := "%s/api/v1/agents/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, id))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result *Agent
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateAgent POST /api/v1/agents
func (c *Client) CreateAgent(agent *Agent) (*Agent, error) {
	m := "POST"
	t := "%s/api/v1/agents"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	payload, err := toJson(agent)
	if err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("agent is nil")
	}

	req, err := http.NewRequest(m, uri, payload)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("create agent response is empty")
	}

	var result Agent
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateRunner POST /api/v1/runners/{name}
func (c *Client) CreateRunner(name string) (*Runner, error) {
	m := "POST"
	t := "%s/api/v1/runners/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, name))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("create runner response is empty")
	}

	var result Runner
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, err
}

// GetRunnerByName GET /api/v1/runners/{name}
func (c *Client) GetRunnerByName(name string) (*Runner, error) {
	m := "GET"
	t := "%s/api/v1/runners/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, name))

	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var result Runner
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, err
}

// UpdateAgent PUT /api/v1/agents/{id}
func (c *Client) UpdateAgent(id string, agent *Agent) (*Agent, error) {
	m := "PUT"
	t := "%s/api/v1/agents/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, id))

	payload, err := toJson(agent)
	if err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("agent is nil")
	}

	req, err := http.NewRequest(m, uri, payload)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create *http.Request")
	}

	body, err := c.doBytesHttpRequest(req)
	if err != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("update agent response is empty")
	}

	var result Agent
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) doBytesHttpRequest(request *http.Request) ([]byte, error) {
	res, err := c.doHttpRequest(request)
	if err != nil {
		return nil, err
	}
	defer closeBody(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil || len(body) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("response body is empty")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) doHttpRequest(request *http.Request) (*http.Response, error) {
	const (
		t = "%s %s"
		a = "ApiKey"
		b = "Bearer"
	)

	header := fmt.Sprintf(t, b, c.userKey)
	if len(c.email) >= 1 && len(c.organization) >= 1 {
		header = fmt.Sprintf(t, a, c.userKey)
	}

	request.Header.Set("Authorization", header)

	return c.client.Do(request)
}

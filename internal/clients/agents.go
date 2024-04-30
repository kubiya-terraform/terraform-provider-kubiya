package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) agents() ([]*Agent, error) {
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

func (c *Client) GetAgents() ([]*Agent, error) {
	agents, err := c.agents()
	return agents, err
}

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

func (c *Client) CreateAgent(agent *Agent) (*Agent, error) {
	m := "POST"
	t := "%s/api/v1/agents"
	uri := c.queryParams(fmt.Sprintf(t, c.host))

	payload, err := toJson(agent)
	if err != nil || payload == nil {
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
	if err != nil || len(body) <= 0 {
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

func (c *Client) GetAgentByName(name string) (*Agent, error) {
	agents, err := c.agents()
	if err != nil {
		return nil, err
	}

	var agent Agent
	for _, item := range agents {
		if item.Name == name {
			agent = Agent{
				Name:                 item.Name,
				Uuid:                 item.Uuid,
				Email:                item.Email,
				Image:                item.Image,
				LlmModel:             item.LlmModel,
				Description:          item.Description,
				Organization:         item.Organization,
				AiInstructions:       item.AiInstructions,
				EnvironmentVariables: make(map[string]string),
				Links:                append(make([]string, 0), item.Links...),
				Owners:               append(make([]string, 0), item.Owners...),
				Runners:              append(make([]string, 0), item.Runners...),
				Secrets:              append(make([]string, 0), item.Secrets...),
				Starters:             append(make([]string, 0), item.Starters...),
				Integrations:         append(make([]string, 0), item.Integrations...),
				AllowedUsers:         append(make([]string, 0), item.AllowedUsers...),
				AllowedGroups:        append(make([]string, 0), item.AllowedGroups...),
			}

			if item.Metadata != nil {
				agent.Metadata = &Metadata{
					CreatedAt:       item.Metadata.CreatedAt,
					LastUpdated:     item.Metadata.LastUpdated,
					UserCreated:     item.Metadata.UserCreated,
					UserLastUpdated: item.Metadata.UserLastUpdated,
				}
			}

			for key, val := range item.EnvironmentVariables {
				agent.EnvironmentVariables[key] = val
			}
			break
		}
	}

	return &agent, err
}

func (c *Client) UpdateAgent(id string, agent *Agent) (*Agent, error) {
	m := "PUT"
	t := "%s/api/v1/agents/%s"
	uri := c.queryParams(fmt.Sprintf(t, c.host, id))

	payload, err := toJson(agent)
	if err != nil || payload == nil {
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
	if err != nil || len(body) <= 0 {
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

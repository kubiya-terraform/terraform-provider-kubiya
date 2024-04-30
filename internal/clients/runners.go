package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) DeleteRunner(name string) error {
	m := "DELETE"
	t := "%s/api/v1/runners/%s-tunnel"
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

func (c *Client) CreateRunner(name, path string) (*Runner, error) {
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
	if err != nil || len(body) <= 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("create runner response is empty")
	}

	var result Runner
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	result.Name = name
	result.Path = toPathYaml(path, name)
	err = c.downloadFile(result.Url, toPathYaml(path, name))

	return &result, err
}

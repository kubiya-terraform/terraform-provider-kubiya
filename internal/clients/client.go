package clients

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
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

func NewClient(k, e, o string) (*Client, error) {
	if len(k) >= 1 {
		host := defaultHost
		timeout := defaultTimeout
		client := &http.Client{Timeout: timeout}
		return &Client{host: host, email: e,
			userKey: k, organization: o, client: client}, nil
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

func (c *Client) downloadFile(uri, path string) error {
	m := "GET"
	req, err := http.NewRequest(m, uri, nil)
	if err != nil || req == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to create *http.Request")
	}

	resp, err := c.doHttpRequest(req)
	if err != nil || resp == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("create runner response is empty")
	}
	defer closeBody(resp.Body)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doBytesHttpRequest(r *http.Request) ([]byte, error) {
	res, err := c.doHttpRequest(r)
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

func (c *Client) doHttpRequest(r *http.Request) (*http.Response, error) {
	const (
		t = "%s %s"
		a = "ApiKey"
		b = "UserKey"
	)

	header := fmt.Sprintf(t, b, c.userKey)
	if len(c.email) >= 1 && len(c.organization) >= 1 {
		header = fmt.Sprintf(t, a, c.userKey)
	}

	r.Header.Set("Authorization", header)

	return c.client.Do(r)
}

func (c *Client) doReaderHttpRequest(r *http.Request) (io.Reader, error) {
	body, err := c.doBytesHttpRequest(r)

	return bytes.NewReader(body), err
}

package clients

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func (c *Client) uri(path string) string {
	const (
		prefix = "/"
		layout = "%s/%s"
	)

	for strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
	}

	return format(layout, c.host, path)
}

func (c *Client) downloadFile(uri, path string) error {
	ctx := context.Background()
	resp, err := c.read(ctx, uri)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) auth(req *http.Request) *http.Request {
	const (
		authLayout = "UserKey %s"
		authHeader = "Authorization"
		source     = "source"
		terraform  = "terraform"
	)

	if req != nil {
		req.Header.Set(authHeader, format(authLayout, c.userKey))
		req.Header.Set(source, terraform)
	}

	return req
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if req != nil {
		req = c.auth(req)
		resp, err := c.client.Do(req)
		if err != nil || resp == nil {
			if err != nil {
				return nil, err
			}

			return nil, eformat("failed to make http request. *http.Response is nil")
		}

		if resp.StatusCode >= http.StatusBadRequest {
			defer closeBody(resp.Body)

			b, e := io.ReadAll(resp.Body)
			if e != nil {
				err = errors.Join(err, e)
			}

			response := string(b)
			statusCode := resp.StatusCode
			requestUrl := req.URL.String()

			err = errors.Join(err, eformat("request :%s has failed. status code: %d, response: %s",
				requestUrl, statusCode, response))
		}

		return resp, err
	}

	return nil, eformat("req of type: *http.Request is nil")
}

func (c *Client) doWithBody(req *http.Request) ([]byte, error) {
	r, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(r.Body)

	return io.ReadAll(r.Body)
}

// createRequest is a helper function to create and configure HTTP requests
func (c *Client) createRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string, qp ...string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}
		return nil, eformat("failed to create *http.Request")
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if len(qp) > 0 {
		req.URL.RawQuery = strings.Join(qp, "&")
	}

	return req.WithContext(ctx), nil
}

// makeRequestWithResponse returns the raw HTTP response
func (c *Client) makeRequestWithResponse(ctx context.Context, method, url string, body io.Reader, headers map[string]string, qp ...string) (*http.Response, error) {
	req, err := c.createRequest(ctx, method, url, body, headers, qp...)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

// makeRequestWithBytes returns the response body as []byte
func (c *Client) makeRequestWithBytes(ctx context.Context, method, url string, body io.Reader, headers map[string]string, qp ...string) ([]byte, error) {
	req, err := c.createRequest(ctx, method, url, body, headers, qp...)
	if err != nil {
		return nil, err
	}

	return c.doWithBody(req)
}

// makeRequest returns the response body as io.Reader
func (c *Client) makeRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string, qp ...string) (io.Reader, error) {
	responseBody, err := c.makeRequestWithBytes(ctx, method, url, body, headers, qp...)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(responseBody), nil
}

// HTTP methods returning io.Reader
func (c *Client) read(ctx context.Context, u string, qp ...string) (io.Reader, error) {
	return c.makeRequest(ctx, http.MethodGet, u, nil, nil, qp...)
}
func (c *Client) readWithJson(ctx context.Context, u string, qp ...string) (io.Reader, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return c.makeRequest(ctx, http.MethodGet, u, nil, headers, qp...)
}

func (c *Client) delete(ctx context.Context, u string, qp ...string) (io.Reader, error) {
	headers := map[string]string{"force": "true"}
	return c.makeRequest(ctx, http.MethodDelete, u, nil, headers, qp...)
}

func (c *Client) deleteWithJson(ctx context.Context, u string, qp ...string) (io.Reader, error) {
	headers := map[string]string{"force": "true", "Content-Type": "application/json"}
	return c.makeRequest(ctx, http.MethodDelete, u, nil, headers, qp...)
}

func (c *Client) create(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	return c.makeRequest(ctx, http.MethodPost, u, b, nil, qp...)
}

func (c *Client) createWithJson(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return c.makeRequest(ctx, http.MethodPost, u, b, headers, qp...)
}

func (c *Client) update(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	return c.makeRequest(ctx, http.MethodPut, u, b, nil, qp...)
}

func (c *Client) updateWithJson(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	return c.makeRequest(ctx, http.MethodPut, u, b, headers, qp...)
}

func (c *Client) readWithBody(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	return c.makeRequest(ctx, http.MethodGet, u, b, nil, qp...)
}

func (c *Client) deleteWithBody(ctx context.Context, u string, b io.Reader, qp ...string) (io.Reader, error) {
	headers := map[string]string{"force": "true"}
	return c.makeRequest(ctx, http.MethodDelete, u, b, headers, qp...)
}

// HTTP methods returning []byte
func (c *Client) readBytes(ctx context.Context, u string, qp ...string) ([]byte, error) {
	return c.makeRequestWithBytes(ctx, http.MethodGet, u, nil, nil, qp...)
}

// HTTP methods returning *http.Response
func (c *Client) deleteResp(ctx context.Context, u string, qp ...string) (*http.Response, error) {
	headers := map[string]string{"force": "true"}
	return c.makeRequestWithResponse(ctx, http.MethodDelete, u, nil, headers, qp...)
}

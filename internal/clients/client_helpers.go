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

func (c *Client) read(ctx context.Context, u string) (io.Reader, error) {
	m := http.MethodGet

	req, err := http.NewRequest(m, u, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

func (c *Client) readBytes(ctx context.Context, u string) ([]byte, error) {
	m := http.MethodGet

	req, err := http.NewRequest(m, u, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	return c.doWithBody(req.WithContext(ctx))
}

func (c *Client) delete(ctx context.Context, u string) (io.Reader, error) {
	m := http.MethodDelete

	req, err := http.NewRequest(m, u, nil)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

func (c *Client) create(ctx context.Context, u string, qp []string, b io.Reader) (io.Reader, error) {
	m := http.MethodPost

	req, err := http.NewRequest(m, u, b)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	if len(qp) > 0 {
		req.URL.RawQuery = strings.Join(qp, "&")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

func (c *Client) update(ctx context.Context, u string, b io.Reader) (io.Reader, error) {
	m := http.MethodPut

	req, err := http.NewRequest(m, u, b)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

func (c *Client) readWithBody(ctx context.Context, u string, b io.Reader) (io.Reader, error) {
	m := http.MethodGet

	req, err := http.NewRequest(m, u, b)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

func (c *Client) deleteWithBody(ctx context.Context, u string, b io.Reader) (io.Reader, error) {
	m := http.MethodDelete

	req, err := http.NewRequest(m, u, b)
	if err != nil || req == nil {
		if err != nil {
			return nil, err
		}

		return nil, eformat("failed to create *http.Request")
	}

	body, err := c.doWithBody(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), err
}

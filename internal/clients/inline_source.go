package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

func newInlineSource(e *entities.InlineSourceModel) (io.Reader, error) {
	type (
		request struct {
			Name      string         `json:"name"`
			Runner    string         `json:"runner"`
			TaskId    string         `json:"task_id"`
			ManagedBy string         `json:"managed_by"`
			Tools     interface{}    `json:"inline_tools"`
			Workflows interface{}    `json:"inline_workflows"`
			Config    map[string]any `json:"dynamic_config"`
		}
	)

	req := &request{
		TaskId:    getTaskId(),
		ManagedBy: getManagedBy(),
		Config:    make(map[string]any),
		Name:      e.Name.ValueString(),
		Runner:    e.Runner.ValueString(),
	}

	if e.Tools.ValueString() != "" && e.Tools.ValueString() != "{}" {
		body := []byte(e.Tools.ValueString())
		if err := json.Unmarshal(body, &req.Tools); err != nil {
			return nil, err
		}
	}

	if e.Config.ValueString() != "" && e.Config.ValueString() != "{}" {
		body := []byte(e.Config.ValueString())
		if err := json.Unmarshal(body, &req.Config); err != nil {
			return nil, err
		}
	}

	if e.Workflows.ValueString() != "" && e.Workflows.ValueString() != "{}" {
		body := []byte(e.Workflows.ValueString())
		if err := json.Unmarshal(body, &req.Workflows); err != nil {
			return nil, err
		}
	}

	return toJson(req)
}

func parseInlineSource(r io.Reader) (*entities.InlineSourceModel, error) {
	type response struct {
		Id             string         `json:"uuid"`
		Type           string         `json:"type"`
		Url            string         `json:"url"`
		Zip            string         `json:"zip"`
		Path           string         `json:"path"`
		Name           string         `json:"name"`
		TaskId         string         `json:"task_id"`
		ManagedBy      string         `json:"managed_by"`
		AgentsCount    int            `json:"connected_agents_count"`
		ToolsCount     int            `json:"connected_tools_count"`
		WorkflowsCount int            `json:"connected_workflows_count"`
		ErrorsCount    int            `json:"errors_count"`
		Config         map[string]any `json:"dynamic_config"`
		Runner         string         `json:"runner"`
	}

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return nil, err
	}

	configData, err := json.Marshal(resp.Config)
	if err != nil {
		return nil, err
	}

	config, err := normalizeJSON(string(configData))
	if err != nil {
		return nil, err
	}

	result := &entities.InlineSourceModel{
		Config: types.StringValue(config),
		Id:     types.StringValue(resp.Id),
		Name:   types.StringValue(resp.Name),
		Type:   types.StringValue(resp.Type),
		Runner: types.StringValue(resp.Runner),
	}

	return result, nil
}

func parseNewInlineSource(r io.Reader) (*entities.InlineSourceModel, error) {
	type response struct {
		Url       string         `json:"url"`
		Type      string         `json:"type"`
		Id        string         `json:"uuid"`
		Name      string         `json:"name"`
		Runner    string         `json:"runner"`
		TaskId    string         `json:"task_id"`
		ManagedBy string         `json:"managed_by"`
		Config    map[string]any `json:"dynamic_config"`
		Errors    []struct {
			File    string `json:"file"`
			Type    string `json:"type"`
			Error   string `json:"error"`
			Details string `json:"details"`
		} `json:"errors,omitempty"`
	}

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return nil, err
	}

	if len(resp.Errors) >= 1 {
		var err error
		const t = "file: %s, type: %s, error: %s, details: %s"
		for _, e := range resp.Errors {
			err = errors.Join(err, eformat(t, e.File, e.Type, e.Error, e.Details))
		}
		return nil, err
	}

	configData, err := json.Marshal(resp.Config)
	if err != nil {
		return nil, err
	}

	config, err := normalizeJSON(string(configData))
	if err != nil {
		return nil, err
	}

	result := &entities.InlineSourceModel{
		Config: types.StringValue(config),
		Id:     types.StringValue(resp.Id),
		Name:   types.StringValue(resp.Name),
		Type:   types.StringValue(resp.Type),
		Runner: types.StringValue(resp.Runner),
	}

	return result, nil
}

func parseInlineSourceTools(r io.Reader, e *entities.InlineSourceModel) error {
	type (
		response struct {
			Id        string           `json:"uuid"`
			Type      string           `json:"type"`
			Tools     interface{}      `json:"tools"`
			Workflows []map[string]any `json:"workflows"`
			Errors    []struct {
				File    string `json:"file"`
				Type    string `json:"type"`
				Error   string `json:"error"`
				Details string `json:"details"`
			} `json:"errors,omitempty"`
		}
	)

	var resp response
	if err := fromJson(r, &resp); err != nil {
		return err
	}

	if len(resp.Errors) >= 1 {
		var err error
		const t = "file: %s, type: %s, error: %s, details: %s"
		for _, respError := range resp.Errors {
			err = errors.Join(err, eformat(t, respError.File,
				respError.Type, respError.Error, respError.Details))
		}
		return err
	}

	if resp.Tools != nil {
		toolsData, err := json.Marshal(resp.Tools)
		if err != nil {
			return err
		}

		normalized, err := normalizeJSON(string(toolsData))
		if err != nil {
			return err
		}

		if normalized == "[]" {
			normalized = ""
		}
		e.Tools = types.StringValue(normalized)
	}

	if len(resp.Workflows) >= 1 {
		workflowList := make([]map[string]any, 0)
		for _, workflow := range resp.Workflows {
			workflowList = append(workflowList, cleanMap(workflow))
		}

		data, err := json.Marshal(workflowList)
		if err != nil {
			return err
		}

		normalized, err := normalizeJSON(string(data))
		if err != nil {
			return err
		}

		if normalized == "[]" {
			normalized = ""
		}
		e.Workflows = types.StringValue(normalized)
	}

	return nil
}

func (c *Client) DeleteInlineSource(ctx context.Context, e *entities.InlineSourceModel) error {
	const (
		requestUri = "/api/v1/sources/%s"
	)
	id := e.Id.ValueString()
	uri := format(requestUri, id)
	resp, err := c.deleteResp(ctx, c.uri(uri))
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return responseBodyError(resp)
	}

	return nil
}

func (c *Client) UpdateInlineSource(ctx context.Context, e *entities.InlineSourceModel) error {
	if e != nil {
		const (
			updateUri   = "/api/v1/sources/%s"
			metadataUri = "/api/v1/sources/%s/metadata"
		)
		id := e.Id.ValueString()
		uri := c.uri(format(updateUri, id))

		body, err := newInlineSource(e)
		if err != nil {
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			return err
		}

		e, err = parseInlineSource(resp)
		if err != nil {
			return err
		}

		uri = c.uri(format(metadataUri, id))

		resp, err = c.read(ctx, uri)
		if err != nil {
			return err
		}

		err = parseInlineSourceTools(resp, e)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("param entity (*entities.InlineSourceModel) is nil")
}

func (c *Client) ReadInlineSource(ctx context.Context, id string) (*entities.InlineSourceModel, error) {
	const (
		readUri     = "/api/v1/sources/%s"
		metadataUri = "/api/v1/sources/%s/metadata"
	)
	uri := format(readUri, id)
	resp, err := c.read(ctx, c.uri(uri))
	if err != nil {
		return nil, err
	}

	result, err := parseInlineSource(resp)
	if err != nil {
		return nil, err
	}

	uri = c.uri(format(metadataUri, id))
	resp, err = c.read(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = parseInlineSourceTools(resp, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) CreateInlineSource(ctx context.Context, e *entities.InlineSourceModel) (*entities.InlineSourceModel, error) {
	if e != nil {
		const (
			createUri   = "/api/v1/sources"
			metadataUri = "/api/v1/sources/%s/metadata"
		)

		uri := c.uri(createUri)

		body, err := newInlineSource(e)
		if err != nil {
			return nil, err
		}

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		result, err := parseNewInlineSource(resp)
		if err != nil {
			return nil, err
		}

		id := result.Id.ValueString()
		uri = c.uri(format(metadataUri, id))

		resp, err = c.read(ctx, uri, "exclude_workflows_tools=true")
		if err != nil {
			return nil, err
		}

		err = parseInlineSourceTools(resp, result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, fmt.Errorf("param entity (*entities.InlineSourceModel) is nil")
}

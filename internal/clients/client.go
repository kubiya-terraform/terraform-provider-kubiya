package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"

	"terraform-provider-kubiya/internal/clients/vendors"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type Client struct {
	host    string
	userKey string
	client  *http.Client
}

func New(key, env string) (*Client, error) {
	// Get logger
	logger := kubiyasentry.GetLogger()

	if len(key) == 0 {
		logger.Error("Failed to create client", "error", "ApiKey is missing or empty")
		return nil, eformat("ApiKey is missing or empty")
	}

	// Create HTTP client with Sentry tracing transport
	client := &http.Client{
		Transport: kubiyasentry.NewHTTPTransport(http.DefaultTransport),
	}

	host := ""
	switch env {
	case "production":
		host = "https://api.kubiya.ai"
	case "staging":
		host = "https://api-staging.dev.kubiya.ai"
	}
	if strings.HasPrefix(env, "http") {
		host = env
	}

	logger.Info("Created Kubiya client",
		"environment", env,
		"host", host,
	)

	kubiyasentry.AddBreadcrumb("client", "Kubiya client created", sentry.LevelInfo, map[string]interface{}{
		"environment": env,
		"host":        host,
	})

	return &Client{userKey: key, client: client, host: host}, nil
}

func (c *Client) self() (*user, error) {
	const (
		path = "/api/v1/users/self"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching current user info", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch current user", "error", err)
		return nil, err
	}

	var result *user
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode user response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched current user")
	return result, nil
}

func (c *Client) state() (*state, error) {
	// Create context with logger
	ctx := context.Background()
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching client state")
	kubiyasentry.AddBreadcrumb("client", "Fetching complete state", sentry.LevelDebug, nil)

	var err error
	var currentState state

	if users, e := c.users(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.userList = append(make([]*user, 0), users...)
	}

	if agents, e := c.agents(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.agentList = append(make([]*agent, 0), agents...)
	}

	if groups, e := c.groups(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.groupList = append(make([]*group, 0), groups...)
	}

	if models, e := c.models(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.modelList = append(make([]string, 0), models...)
	}

	if runners, e := c.runners(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.runnerList = append(make([]*runner, 0), runners...)
	}

	if secrets, e := c.secrets(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.secretList = append(make([]*secret, 0), secrets...)
	}

	if sources, e := c.sources(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.sourceList = append(make([]*source, 0), sources...)
	}

	if webhooks, e := c.webhooks(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.webhookList = append(make([]*webhook, 0), webhooks...)
	}

	if integrations, e := c.integrations(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.integrationList = append(make([]*integration, 0), integrations...)
	}

	if knowledgeList, e := c.knowledge(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.knowledgeList = append(make([]*knowledge, 0), knowledgeList...)
	}

	if externalKnowledgeList, e := c.externalKnowledge(); e != nil {
		err = errors.Join(err, e)
	} else {
		currentState.externalKnowledgeList = append(make([]*vendors.BaseExternalKnowledge, 0), externalKnowledgeList...)
	}

	if err != nil {
		logger.Error("Failed to fetch complete state", "error", err)
	} else {
		logger.Debug("Successfully fetched client state",
			"users", len(currentState.userList),
			"agents", len(currentState.agentList),
			"groups", len(currentState.groupList),
			"runners", len(currentState.runnerList),
			"secrets", len(currentState.secretList),
			"sources", len(currentState.sourceList),
			"webhooks", len(currentState.webhookList),
			"integrations", len(currentState.integrationList),
			"knowledge", len(currentState.knowledgeList),
			"external_knowledge", len(currentState.externalKnowledgeList),
		)
	}

	return &currentState, err
}

func (c *Client) users() ([]*user, error) {
	const (
		path = "/api/v1/users"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching users list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch users", "error", err)
		return nil, err
	}

	var result []*user
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode users response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched users", "count", len(result))
	return result, nil
}

func (c *Client) agents() ([]*agent, error) {
	const (
		path = "/api/v1/agents"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching agents list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch agents", "error", err)
		return nil, err
	}

	var result []*agent
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode agents response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched agents", "count", len(result))
	return result, nil
}

func (c *Client) groups() ([]*group, error) {
	const (
		path = "/api/v1/manage/groups"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching groups list", "path", path)

	resp, err := c.readBytes(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch groups", "error", err)
		return nil, err
	}

	var result []*group
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode groups response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched groups", "count", len(result))
	return result, nil
}

func (c *Client) models() ([]string, error) {
	const (
		sep  = ","
		path = "/api/v1/featureflags"
		body = `["supported_llm_models"]`
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching supported models", "path", path)

	payload := strings.NewReader(body)

	resp, err := c.create(ctx, uri, payload)
	if err != nil {
		logger.Error("Failed to fetch models", "error", err)
		return nil, err
	}

	tmp := &struct {
		Models string `json:"supported_llm_models"`
	}{}

	if err = json.NewDecoder(resp).Decode(tmp); err != nil {
		logger.Error("Failed to decode models response", "error", err)
		return nil, err
	}

	var result []string

	for _, item := range strings.Split(tmp.Models, sep) {
		result = append(result, strings.TrimSpace(item))
	}

	logger.Debug("Successfully fetched models", "count", len(result))
	return result, nil
}

func (c *Client) runners() ([]*runner, error) {
	const (
		path = "/api/v3/runners"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching runners list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch runners", "error", err)
		return nil, err
	}

	var result []*runner
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode runners response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched runners", "count", len(result))
	return result, nil
}

func (c *Client) secrets() ([]*secret, error) {
	const (
		path = "/api/v2/secrets"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching secrets list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch secrets", "error", err)
		return nil, err
	}

	var result []*secret
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode secrets response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched secrets", "count", len(result))
	return result, nil
}

func (c *Client) sources() ([]*source, error) {
	const (
		path = "/api/v1/sources"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching sources list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch sources", "error", err)
		return nil, err
	}

	result, err := newSources(resp)
	if err != nil {
		logger.Error("Failed to parse sources response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched sources", "count", len(result))
	return result, nil
}

func (c *Client) webhooks() ([]*webhook, error) {
	const (
		path = "/api/v1/event"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching webhooks list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch webhooks", "error", err)
		return nil, err
	}

	var result []*webhook
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode webhooks response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched webhooks", "count", len(result))
	return result, nil
}

func (c *Client) knowledge() ([]*knowledge, error) {
	const (
		path = "/api/v1/knowledge"
	)

	uri := c.uri(path)
	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching knowledge list", "path", path)

	resp, err := c.read(ctx, uri)
	if err != nil {
		logger.Error("Failed to fetch knowledge", "error", err)
		return nil, err
	}

	var result []*knowledge
	err = json.NewDecoder(resp).Decode(&result)
	if err != nil {
		logger.Error("Failed to decode knowledge response", "error", err)
		return nil, err
	}

	logger.Debug("Successfully fetched knowledge", "count", len(result))
	return result, nil
}

func (c *Client) integrations() ([]*integration, error) {
	const (
		pathIntegration = "/api/v2/integrations"
	)

	ctx := context.Background()

	// Ensure logger is in context
	logger := kubiyasentry.GetLogger()
	ctx = kubiyasentry.ContextWithLogger(ctx, logger)

	logger.Debug("Fetching integrations list", "path", pathIntegration)

	result := []*integration{
		{Name: "slack"},
		{Name: "kubernetes"},
	}

	// Only call the integrations endpoint
	resp, err := c.read(ctx, c.uri(pathIntegration))
	if err != nil {
		logger.Error("Failed to fetch integrations", "error", err)
		return nil, err
	}

	var tmpList []*integrationApi
	err = json.NewDecoder(resp).Decode(&tmpList)
	if err != nil {
		logger.Error("Failed to decode integrations response", "error", err)
		return nil, err
	}

	for _, item := range tmpList {
		result = append(result, &integration{
			Name: item.Name,
		})
	}

	logger.Debug("Successfully fetched integrations", "count", len(result))
	return result, nil
}

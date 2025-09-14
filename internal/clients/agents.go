package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

type task struct {
	Name        string `json:"name"`
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
}

type agent struct {
	Name      string `json:"name"`
	Uuid      string `json:"uuid"`
	TaskId    string `json:"task_id"`
	ManagedBy string `json:"managed_by"`
	Email     string `json:"email,omitempty"`
	Image     string `json:"image,omitempty"`

	Links        []string  `json:"links"`
	Tools        []string  `json:"tools"`
	Tasks        []task    `json:"tasks"`
	Sources      []string  `json:"sources"`
	Secrets      []string  `json:"secrets"`
	Starters     []starter `json:"starters"`
	Integrations []string  `json:"integrations"`
	Users        []string  `json:"allowed_users"`
	Groups       []string  `json:"allowed_groups"`
	Owners       []string  `json:"owners,omitempty"`
	Runners      []string  `json:"runners,omitempty"`
	IsDebugMode  bool      `json:"is_debug_mode,omitempty"`

	Metadata  *metadata         `json:"metadata"`
	Variables map[string]string `json:"environment_variables"`
	MCPServer *mcpServer        `json:"mcp_server,omitempty"`

	LlmModel       string `json:"llm_model,omitempty"`
	Description    string `json:"description,omitempty"`
	Organization   string `json:"organization,omitempty"`
	AiInstructions string `json:"ai_instructions,omitempty"`
}

type starter struct {
	Command string `json:"command"`
	Name    string `json:"display_name"`
}

type mcpServer struct {
	Type    string            `json:"type"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type metadata struct {
	CreatedAt       string `json:"created_at"`
	LastUpdated     string `json:"last_updated"`
	UserCreated     string `json:"user_created"`
	UserLastUpdated string `json:"user_last_updated"`
}

func toAgent(a *entities.AgentModel, cs *state) (*agent, error) {
	var err error
	var validRunner bool

	result := &agent{
		Uuid:           a.Id.ValueString(),
		Name:           a.Name.ValueString(),
		Image:          a.Image.ValueString(),
		LlmModel:       a.Model.ValueString(),
		Description:    a.Description.ValueString(),
		AiInstructions: a.Instructions.ValueString(),
		Runners:        []string{a.Runner.ValueString()},

		Email:        "",
		Organization: "",
		IsDebugMode:  a.IsDebugMode.ValueBool(),

		Owners: make([]string, 0),

		Links:        make([]string, 0),
		Tools:        make([]string, 0),
		Users:        make([]string, 0),
		Groups:       make([]string, 0),
		Sources:      make([]string, 0),
		Secrets:      make([]string, 0),
		Integrations: make([]string, 0),
		Variables:    make(map[string]string),

		Tasks:    make([]task, 0),
		Starters: make([]starter, 0),
	}

	for _, v := range cs.runnerList {
		item := a.Runner.ValueString()
		if validRunner = equal(v.Name, item); validRunner {
			break
		}
	}

	if !validRunner {
		item := a.Runner
		err = errors.Join(err, eformat("runner \"%s\" don't exist", item))
	}

	if len(a.Tools.Elements()) >= 6 {
		err = errors.Join(err, eformat("tools field can have no more than 5 elements"))
	}

	for _, v := range a.Tasks {
		result.Tasks = append(result.Tasks, task{
			Name:        v.Name,
			Prompt:      v.Prompt,
			Description: v.Description,
		})
	}

	for _, v := range a.Starters {
		result.Starters = append(result.Starters, starter{
			Name:    v.Name,
			Command: v.Command,
		})
	}

	for _, v := range a.Links.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			str := v.String()
			result.Links = append(result.Links, strings.ReplaceAll(str, "\"", ""))
		}
	}

	for _, v := range a.Tools.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			str := v.String()
			result.Tools = append(result.Tools, strings.ReplaceAll(str, "\"", ""))
		}
	}

	for _, v := range a.Users.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.userList {
				if found = equal(i.Name, item) ||
					equal(i.Email, item); found {
					result.Users = append(result.Users, i.UUID)
					break
				}
			}
			if !found {
				err = errors.Join(err, fmt.Errorf("user \"%s\" don't exist", item))
			}
		}
	}

	for _, v := range a.Groups.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.groupList {
				if found = equal(i.Name, item); found {
					result.Groups = append(result.Groups, i.UUID)
					break
				}
			}
			if !found {
				err = errors.Join(err, fmt.Errorf("group \"%s\" don't exist", v))
			}
		}
	}

	for _, v := range a.Secrets.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.secretList {
				if found = equal(i.Name, item); found {
					result.Secrets = append(result.Secrets, i.Name)
					break
				}
			}
			if !found {
				err = errors.Join(err, fmt.Errorf("secret \"%s\" don't exist", v))
			}
		}
	}

	for _, v := range a.Sources.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.sourceList {
				if found = equal(i.Name, item); found {
					result.Sources = append(result.Sources, i.Id)
					break
				}
			}
			if !found {
				err = errors.Join(err, fmt.Errorf("source \"%s\" don't exist", v))
			}
		}
	}

	for _, v := range a.Integrations.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.integrationList {
				if found = equal(i.Name, item); found {
					result.Integrations = append(result.Integrations, i.Name)
					break
				}
			}

			if !found {
				err = errors.Join(err, fmt.Errorf("integration \"%s\" don't exist", v))
			}
		}
	}

	for key, value := range a.Variables.Elements() {
		result.Variables[key] = strings.ReplaceAll(value.String(), "\"", "")
	}

	// Convert MCP server configuration
	if a.MCPServer != nil && !a.MCPServer.Type.IsNull() && !a.MCPServer.Type.IsUnknown() {
		result.MCPServer = &mcpServer{
			Type: a.MCPServer.Type.ValueString(),
		}

		// Handle STDIO-specific fields
		if !a.MCPServer.Command.IsNull() && !a.MCPServer.Command.IsUnknown() {
			result.MCPServer.Command = a.MCPServer.Command.ValueString()
		}

		if !a.MCPServer.Args.IsNull() && !a.MCPServer.Args.IsUnknown() {
			result.MCPServer.Args = make([]string, 0)
			for _, v := range a.MCPServer.Args.Elements() {
				if !v.IsNull() && !v.IsUnknown() {
					str := v.String()
					result.MCPServer.Args = append(result.MCPServer.Args,
						strings.ReplaceAll(str, "\"", ""))
				}
			}
		}

		if !a.MCPServer.Env.IsNull() && !a.MCPServer.Env.IsUnknown() {
			result.MCPServer.Env = make(map[string]string)
			for key, value := range a.MCPServer.Env.Elements() {
				result.MCPServer.Env[key] = strings.ReplaceAll(value.String(), "\"", "")
			}
		}

		// Handle SSE-specific fields
		if !a.MCPServer.URL.IsNull() && !a.MCPServer.URL.IsUnknown() {
			result.MCPServer.URL = a.MCPServer.URL.ValueString()
		}

		if !a.MCPServer.Headers.IsNull() && !a.MCPServer.Headers.IsUnknown() {
			result.MCPServer.Headers = make(map[string]string)
			for key, value := range a.MCPServer.Headers.Elements() {
				result.MCPServer.Headers[key] = strings.ReplaceAll(value.String(), "\"", "")
			}
		}
	}

	if valid := slices.Contains(cs.modelList, result.LlmModel); !valid {
		model := result.LlmModel
		models := strings.Join(cs.modelList, ",")
		err = errors.Join(err, eformat("LLM Model \"%s\" not valid. [%s]", model, models))
	}

	return result, err
}

func fromAgent(a *agent, cs *state) (*entities.AgentModel, error) {
	var err error
	result := &entities.AgentModel{
		Id:           types.StringValue(a.Uuid),
		Name:         types.StringValue(a.Name),
		Image:        types.StringValue(a.Image),
		Model:        types.StringValue(a.LlmModel),
		IsDebugMode:  types.BoolValue(a.IsDebugMode),
		Description:  types.StringValue(a.Description),
		Instructions: types.StringValue(a.AiInstructions),
	}

	usersList := make([]string, 0)
	groupList := make([]string, 0)
	sourceList := make([]string, 0)

	if a.Metadata != nil {
		for _, u := range cs.userList {
			if equal(u.UUID, a.Metadata.UserCreated) {
				result.Owner = types.StringValue(u.Email)
				break
			}
		}

		result.CreatedAt = types.StringValue(a.Metadata.UserCreated)
	}

	if len(a.Runners) >= 1 {
		result.Runner = types.StringValue(a.Runners[0])
	}

	if len(a.Tasks) >= 1 {
		result.Tasks = make([]entities.TaskModel, 0)
		for _, t := range a.Tasks {
			result.Tasks = append(result.Tasks, entities.TaskModel{
				Name:        t.Name,
				Prompt:      t.Prompt,
				Description: t.Description,
			})
		}
	}

	if len(a.Starters) >= 1 {
		result.Starters = make([]entities.StarterModel, 0)
		for _, t := range a.Starters {
			result.Starters = append(result.Starters, entities.StarterModel{
				Name:    t.Name,
				Command: t.Command,
			})
		}
	}

	for _, t := range a.Users {
		for _, u := range cs.userList {
			if equal(u.UUID, t) {
				usersList = append(usersList, u.Email)
				break
			}
		}
	}

	for _, t := range a.Groups {
		for _, g := range cs.groupList {
			if equal(g.UUID, t) {
				groupList = append(groupList, g.Name)
				break
			}
		}
	}

	for _, t := range a.Sources {
		for _, s := range cs.sourceList {
			if equal(s.Id, t) {
				sourceList = append(sourceList, s.Name)
				break
			}
		}
	}

	result.Tools = toListStringType(a.Tools, err)

	result.Links = toListStringType(a.Links, err)

	result.Variables = toMapType(a.Variables, err)

	result.Users = toListStringType(usersList, err)

	result.Groups = toListStringType(groupList, err)

	result.Secrets = toListStringType(a.Secrets, err)

	result.Sources = toListStringType(sourceList, err)

	result.Integrations = toListStringType(a.Integrations, err)

	// Convert MCP server configuration
	if a.MCPServer != nil {
		result.MCPServer = &entities.MCPServerModel{
			Type: types.StringValue(a.MCPServer.Type),
		}

		// Handle STDIO-specific fields
		if a.MCPServer.Command != "" {
			result.MCPServer.Command = types.StringValue(a.MCPServer.Command)
		} else {
			result.MCPServer.Command = types.StringNull()
		}

		if len(a.MCPServer.Args) > 0 {
			result.MCPServer.Args = toListStringType(a.MCPServer.Args, err)
		} else {
			result.MCPServer.Args = types.ListNull(types.StringType)
		}

		if len(a.MCPServer.Env) > 0 {
			result.MCPServer.Env = toMapType(a.MCPServer.Env, err)
		} else {
			result.MCPServer.Env = types.MapNull(types.StringType)
		}

		// Handle SSE-specific fields
		if a.MCPServer.URL != "" {
			result.MCPServer.URL = types.StringValue(a.MCPServer.URL)
		} else {
			result.MCPServer.URL = types.StringNull()
		}

		if len(a.MCPServer.Headers) > 0 {
			result.MCPServer.Headers = toMapType(a.MCPServer.Headers, err)
		} else {
			result.MCPServer.Headers = types.MapNull(types.StringType)
		}
	}

	return result, err
}

func (c *Client) DeleteAgent(ctx context.Context, e *entities.AgentModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/agents/%s", id)

		// Log deletion attempt
		logger.Info("Deleting agent",
			"agent_id", id,
			"path", path,
		)

		// Add breadcrumb
		kubiyasentry.AddBreadcrumb("client", "Deleting agent via API", sentry.LevelInfo, map[string]interface{}{
			"method":   "DeleteAgent",
			"agent_id": id,
			"uri":      path,
		})

		_, err := c.delete(ctx, c.uri(path))
		if err != nil {
			logger.Error("Failed to delete agent",
				"agent_id", id,
				"error", err,
			)
			return err
		}

		logger.Info("Successfully deleted agent",
			"agent_id", id,
		)
		return nil
	}

	logger.Error("DeleteAgent called with nil entity")
	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) UpdateAgent(ctx context.Context, e *entities.AgentModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	if e != nil {
		id := e.Id.ValueString()
		name := e.Name.ValueString()

		// Log update attempt
		logger.Info("Updating agent",
			"agent_id", id,
			"agent_name", name,
		)

		// Add breadcrumb
		kubiyasentry.AddBreadcrumb("client", "Updating agent via API", sentry.LevelInfo, map[string]interface{}{
			"method":     "UpdateAgent",
			"agent_id":   id,
			"agent_name": name,
		})

		cs, err := c.state()
		if err != nil {
			logger.Error("Failed to get state for agent update",
				"agent_id", id,
				"error", err,
			)
			return err
		}

		e.Owner = types.StringNull()
		uri := c.uri(format("/api/v1/agents/%s", id))

		data, err := toAgent(e, cs)
		if err != nil {
			return err
		}

		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to update agent",
				"agent_id", id,
				"error", err,
			)
			return err
		}

		var r *agent
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode agent update response",
				"agent_id", id,
				"error", err,
			)
			return err
		}

		e, err = fromAgent(r, cs)
		if err != nil {
			logger.Error("Failed to convert updated agent from API response",
				"agent_id", id,
				"error", err,
			)
			return err
		}

		logger.Info("Successfully updated agent",
			"agent_id", id,
			"agent_name", r.Name,
		)
		return nil
	}

	logger.Error("UpdateAgent called with nil entity")
	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) ReadAgent(ctx context.Context, id string) (*entities.AgentModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Log read attempt
	logger.Debug("Reading agent",
		"agent_id", id,
	)

	// Add breadcrumb
	kubiyasentry.AddBreadcrumb("client", "Reading agent via API", sentry.LevelDebug, map[string]interface{}{
		"method":   "ReadAgent",
		"agent_id": id,
	})

	cs, err := c.state()
	if err != nil {
		logger.Error("Failed to get state for agent read",
			"agent_id", id,
			"error", err,
		)
		return nil, err
	}

	path := format("/api/v1/agents/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to read agent",
			"agent_id", id,
			"error", err,
		)
		return nil, err
	}

	var r *agent
	err = json.NewDecoder(resp).Decode(&r)
	if err != nil {
		logger.Error("Failed to decode agent read response",
			"agent_id", id,
			"error", err,
		)
		return nil, err
	}

	entity, err := fromAgent(r, cs)
	if err != nil || entity == nil {
		if err != nil {
			logger.Error("Failed to convert agent from API response",
				"agent_id", id,
				"error", err,
			)
			return nil, err
		}
		logger.Warn("Agent not found",
			"agent_id", id,
		)
		return nil, eformat("Agent %s not found", id)
	}

	logger.Debug("Successfully read agent",
		"agent_id", id,
		"agent_name", entity.Name.ValueString(),
	)
	return entity, nil
}

func (c *Client) CreateAgent(ctx context.Context, e *entities.AgentModel) (*entities.AgentModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateAgent")
		span.SetData("client.resource_type", "agent")
	}

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating agent via API", sentry.LevelInfo, map[string]interface{}{
		"method": "CreateAgent",
		"uri":    "/api/v1/agents",
	})

	if e != nil {
		name := e.Name.ValueString()

		// Log creation attempt
		logger.Info("Creating agent",
			"agent_name", name,
			"model", e.Model.ValueString(),
			"runner", e.Runner.ValueString(),
		)

		cs, err := c.state()
		if err != nil {
			logger.Error("Failed to get state for agent creation",
				"agent_name", name,
				"error", err,
			)
			return nil, err
		}

		data, err := toAgent(e, cs)
		if err != nil {
			logger.Error("Failed to convert agent to API format",
				"agent_name", name,
				"error", err,
			)
			return nil, err
		}

		data.ManagedBy, data.TaskId = managedBy()

		body, err := toJson(data)
		if err != nil {
			logger.Error("Failed to marshal agent to JSON",
				"agent_name", e.Name.ValueString(),
				"error", err,
			)
			return nil, err
		}

		uri := c.uri("/api/v1/agents")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			logger.Error("Failed to create agent",
				"agent_name", e.Name.ValueString(),
				"error", err,
			)
			return nil, err
		}

		var r *agent
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			logger.Error("Failed to decode agent creation response",
				"agent_name", e.Name.ValueString(),
				"error", err,
			)
			return nil, err
		}

		entity, err := fromAgent(r, cs)
		if err != nil {
			logger.Error("Failed to convert created agent from API response",
				"agent_name", e.Name.ValueString(),
				"error", err,
			)
			return nil, err
		}

		logger.Info("Successfully created agent",
			"agent_id", r.Uuid,
			"agent_name", r.Name,
		)
		return entity, nil
	}

	logger.Error("CreateAgent called with nil entity")
	return e, fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

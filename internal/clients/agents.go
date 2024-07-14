package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type task struct {
	Name        string `json:"name"`
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
}

type agent struct {
	Name  string `json:"name"`
	Uuid  string `json:"uuid"`
	Email string `json:"email,omitempty"`
	Image string `json:"image,omitempty"`

	Links        []string  `json:"links"`
	Tools        []string  `json:"tools"`
	Tasks        []task    `json:"tasks"`
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

	LlmModel       string `json:"llm_model,omitempty"`
	Description    string `json:"description,omitempty"`
	Organization   string `json:"organization,omitempty"`
	AiInstructions string `json:"ai_instructions,omitempty"`
}

type starter struct {
	Command string `json:"command"`
	Name    string `json:"display_name"`
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

	result.Tools = toListStringType(a.Tools, err)

	result.Links = toListStringType(a.Links, err)

	result.Variables = toMapType(a.Variables, err)

	result.Users = toListStringType(usersList, err)

	result.Groups = toListStringType(groupList, err)

	result.Secrets = toListStringType(a.Secrets, err)

	result.Integrations = toListStringType(a.Integrations, err)

	return result, err
}

func (c *Client) ReadAgent(_ context.Context, e *entities.AgentModel) error {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := e.Id
		name := e.Name

		for _, a := range cs.agentList {
			if equal(a.Uuid, id.ValueString()) ||
				equal(a.Name, name.ValueString()) {
				e, err = fromAgent(a, cs)
				break
			}
		}

		return err
	}

	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) DeleteAgent(ctx context.Context, e *entities.AgentModel) error {
	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/agents/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		return err
	}

	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) UpdateAgent(ctx context.Context, e *entities.AgentModel) error {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := e.Id.ValueString()
		e.Owner = types.StringNull()
		uri := c.uri(format("/api/v1/agents/%s", id))

		data, err := toAgent(e, cs)
		if err != nil {
			return err
		}

		body, err := toJson(data)
		if err != nil {
			return err
		}

		resp, err := c.update(ctx, uri, body)
		if err != nil {
			return err
		}

		var r *agent
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return err
		}

		e, err = fromAgent(r, cs)
		return err
	}
	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) CreateAgent(ctx context.Context, e *entities.AgentModel) (*entities.AgentModel, error) {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		data, err := toAgent(e, cs)
		if err != nil {
			return nil, err
		}

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/agents")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		var r *agent
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		return fromAgent(r, cs)
	}

	return e, fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

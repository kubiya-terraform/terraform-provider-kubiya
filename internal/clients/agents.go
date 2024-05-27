package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	Tasks        []task    `json:"tasks"`
	Secrets      []string  `json:"secrets"`
	Starters     []starter `json:"starters"`
	Integrations []string  `json:"integrations"`
	Users        []string  `json:"allowed_users"`
	Groups       []string  `json:"allowed_groups"`
	Owners       []string  `json:"owners,omitempty"`
	Runners      []string  `json:"runners,omitempty"`

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

		Owners: make([]string, 0),

		Links:        make([]string, 0),
		Users:        make([]string, 0),
		Groups:       make([]string, 0),
		Secrets:      make([]string, 0),
		Integrations: make([]string, 0),
		Variables:    make(map[string]string),

		Tasks:    make([]task, 0),
		Starters: make([]starter, 0),
	}

	for _, v := range cs.runners {
		item := a.Runner.ValueString()
		if validRunner = equal(v.Name, item); validRunner {
			break
		}
	}

	if !validRunner {
		item := a.Runner
		err = errors.Join(err, fmt.Errorf("runner \"%s\" don't exist", item))
	}

	for _, v := range a.Tasks {
		result.Tasks = append(result.Tasks, task{
			Name:        v.Name,
			Prompt:      v.Prompt,
			Description: v.Description,
		})
	}

	for _, v := range a.Links.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			str := v.String()
			result.Links = append(result.Links, strings.ReplaceAll(str, "\"", ""))
		}
	}

	for _, v := range a.Users.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.users {
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
			for _, i := range cs.groups {
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
			for _, i := range cs.secrets {
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

	for _, v := range a.Starters {
		result.Starters = append(result.Starters, starter{
			Name:    v.Name,
			Command: v.Command,
		})
	}

	for _, v := range a.Integrations.Elements() {
		if !v.IsNull() && !v.IsUnknown() {
			found := false
			str := v.String()
			item := strings.ReplaceAll(str, "\"", "")
			for _, i := range cs.integrations {
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

	return result, err
}

func fromAgent(a *agent, cs *state) (*entities.AgentModel, error) {
	var err error
	result := &entities.AgentModel{
		Id:           types.StringValue(a.Uuid),
		Name:         types.StringValue(a.Name),
		Image:        types.StringValue(a.Image),
		Model:        types.StringValue(a.LlmModel),
		Description:  types.StringValue(a.Description),
		Instructions: types.StringValue(a.AiInstructions),

		//Tasks:    make([]entities.TaskModel, 0),
		//Starters: make([]entities.StarterModel, 0),
	}

	usersList := make([]string, 0)
	groupList := make([]string, 0)

	if a.Metadata != nil {
		for _, u := range cs.users {
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
		result.Tasks = make([]entities.TaskModel, 0)
		for _, t := range a.Starters {
			result.Starters = append(result.Starters, entities.StarterModel{
				Name:    t.Name,
				Command: t.Command,
			})
		}
	}

	for _, t := range a.Users {
		for _, u := range cs.users {
			if equal(u.UUID, t) {
				usersList = append(usersList, u.Email)
				break
			}
		}
	}

	for _, t := range a.Groups {
		for _, g := range cs.groups {
			if equal(g.UUID, t) {
				groupList = append(groupList, g.Name)
				break
			}
		}
	}

	//for _, t := range a.Starters {
	//	result.Starters = make([]entities.StarterModel, 0)
	//	result.Starters = append(result.Starters, entities.StarterModel{
	//		Name:    t.Name,
	//		Command: t.Command,
	//	})
	//}

	//if m := toMapType(a.Variables, err); err == nil {
	//	result.Variables = m
	//}

	result.Links = toListStringType(a.Links, err)

	result.Variables = toMapType(a.Variables, err)

	result.Users = toListStringType(usersList, err)

	result.Groups = toListStringType(groupList, err)

	result.Secrets = toListStringType(a.Secrets, err)

	result.Integrations = toListStringType(a.Integrations, err)

	//if len(a.Tasks) >= 1 {
	//	result.Tasks = make([]entities.TaskModel, 0)
	//	for _, item := range a.Tasks {
	//		result.Tasks = append(result.Tasks, entities.TaskModel{
	//			Name:        item.Name,
	//			Prompt:      item.Prompt,
	//			Description: item.Description,
	//		})
	//	}
	//}
	//
	//if len(a.Links) >= 1 {
	//	result.Links = make([]string, 0)
	//	result.Links = append(result.Links, item)
	//}
	//
	//if len(a.Users) >= 1 {
	//	result.Users = make([]string, 0)
	//	for _, item := range a.Users {
	//		for _, u := range cs.users {
	//			if equal(u.UUID, item) {
	//				result.Users = append(result.Users, u.Email)
	//				break
	//			}
	//		}
	//	}
	//}
	//
	//if len(a.Groups) >= 1 {
	//	result.Groups = make([]string, 0)
	//
	//	for _, item := range a.Groups {
	//		for _, g := range cs.groups {
	//			if equal(g.UUID, item) {
	//				result.Groups = append(result.Groups, g.Name)
	//				break
	//			}
	//		}
	//	}
	//}
	//
	//if len(a.Secrets) >= 1 {
	//	result.Secrets = make([]string, 0)
	//	for _, item := range a.Secrets {
	//		result.Secrets = append(result.Secrets, item)
	//	}
	//}
	//
	//if len(a.Starters) >= 1 {
	//	result.Starters = make([]entities.StarterModel, 0)
	//	for _, item := range a.Starters {
	//		result.Starters = append(result.Starters, entities.StarterModel{
	//			Name:    item.Name,
	//			Command: item.Command,
	//		})
	//	}
	//}
	//
	//if len(a.Variables) >= 1 {
	//	variables := map[string]string{}
	//	result.Variables = &map[string]string{}
	//	for key, value := range a.Variables {
	//		variables[key] = value
	//		//result.Variables[key] = value
	//	}
	//
	//	result.Variables = &variables
	//}
	//
	//if len(a.Integrations) >= 1 {
	//	result.Integrations = make([]string, 0)
	//	for _, item := range a.Integrations {
	//		result.Integrations = append(result.Integrations, item)
	//	}
	//}

	return result, err
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

func (c *Client) ReadAgent(_ context.Context, e *entities.AgentModel) (*entities.AgentModel, error) {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		id := e.Id
		name := e.Name

		var entity *entities.AgentModel
		for _, a := range cs.agents {
			if equal(a.Uuid, id.ValueString()) ||
				equal(a.Name, name.ValueString()) {
				entity, err = fromAgent(a, cs)
				break
			}
		}

		return entity, err
	}

	return e, fmt.Errorf("param entity (*entities.AgentModel) is nil")
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

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

type agent struct {
	Name                 string            `json:"name"`
	Uuid                 string            `json:"uuid"`
	Email                string            `json:"email,omitempty"`
	Image                string            `json:"image,omitempty"`
	Links                []string          `json:"links"`
	Owners               []string          `json:"owners,omitempty"`
	Runners              []string          `json:"runners,omitempty"`
	Secrets              []string          `json:"secrets"`
	Starters             []Starter         `json:"starters"`
	Metadata             *metadata         `json:"metadata"`
	LlmModel             string            `json:"llm_model,omitempty"`
	Description          string            `json:"description,omitempty"`
	Integrations         []string          `json:"integrations"`
	Organization         string            `json:"organization,omitempty"`
	AllowedUsers         []string          `json:"allowed_users"`
	AllowedGroups        []string          `json:"allowed_groups"`
	AiInstructions       string            `json:"ai_instructions,omitempty"`
	EnvironmentVariables map[string]string `json:"environment_variables"`
}

type metadata struct {
	CreatedAt       string `json:"created_at"`
	LastUpdated     string `json:"last_updated"`
	UserCreated     string `json:"user_created"`
	UserLastUpdated string `json:"user_last_updated"`
}

type Starter struct {
	DisplayName string `json:"display_name"`
	Command     string `json:"command"`
}

func toAgent(a *entities.AgentModel, cs *state) (*agent, error) {
	var err error
	users := make([]string, 0)
	groups := make([]string, 0)
	runners := make([]string, 0)
	secrets := make([]string, 0)
	envs := a.Variables.ValueString()
	integrations := make([]string, 0)
	envVariables := make(map[string]string)

	for _, obj := range strings.Split(envs, ",") {
		if kv := strings.Split(obj, ":"); len(kv) == 2 {
			key := kv[0]
			value := kv[1]

			key = strings.TrimSpace(key)
			key = strings.ReplaceAll(key, "{", "")
			key = strings.ReplaceAll(key, "}", "")
			key = strings.ReplaceAll(key, "\"", "")

			value = strings.ReplaceAll(value, "{", "")
			value = strings.ReplaceAll(value, "}", "")

			envVariables[key] = strings.ReplaceAll(value, "\"", "")
		}
	}

	for _, item := range stringList(a.Users.ValueString()) {
		if len(item) <= 0 {
			continue
		}
		found := false
		for _, i := range cs.users {
			if found = equal(i.Name, item) ||
				equal(i.Email, item); found {
				users = append(users, i.UUID)
				break
			}
		}
		if !found {
			err = errors.Join(err, fmt.Errorf("User \"%s\" don't exist", item))
		}
	}

	for _, item := range stringList(a.Groups.ValueString()) {
		if len(item) <= 0 {
			continue
		}
		found := false
		for _, i := range cs.groups {
			if found = equal(i.Name, item); found {
				groups = append(groups, i.UUID)
				break
			}
		}
		if !found {
			err = errors.Join(err, fmt.Errorf("Group \"%s\" don't exist", item))
		}
	}

	for _, item := range stringList(a.Runners.ValueString()) {
		if len(item) <= 0 {
			continue
		}
		found := false
		for _, i := range cs.runners {
			if found = equal(i.Name, item); found {
				runners = append(runners, i.Name)
				break
			}
		}
		if !found {
			err = errors.Join(err, fmt.Errorf("Runner \"%s\" don't exist", item))
		}
	}

	for _, item := range stringList(a.Secrets.ValueString()) {
		if len(item) <= 0 {
			continue
		}
		found := false
		for _, i := range cs.secrets {
			if found = equal(i.Name, item); found {
				secrets = append(secrets, i.Name)
				break
			}
		}
		if !found {
			err = errors.Join(err, fmt.Errorf("Secret \"%s\" don't exist", item))
		}
	}

	for _, item := range stringList(a.Integrations.ValueString()) {
		if len(item) <= 0 {
			continue
		}
		found := false
		for _, i := range cs.integrations {
			if found = equal(i.Name, item); found {
				integrations = append(integrations, i.Name)
				break
			}
		}
		if !found {
			err = errors.Join(err, fmt.Errorf("Integration \"%s\" don't exist", item))
		}
	}

	if len(runners) <= 0 {
		err = errors.Join(err, fmt.Errorf("runners cannot be empty. you must have at least one"))
	}

	var starters []Starter
	if a.Starters.ValueString() == "" {
		starters = []Starter{}
	} else {
		input := a.Starters.ValueString()
		err = json.Unmarshal([]byte(input), &starters)
		if err != nil {
			return nil, err
		}
	}

	return &agent{
		Uuid:           a.Id.ValueString(),
		Name:           a.Name.ValueString(),
		Image:          a.Image.ValueString(),
		LlmModel:       a.Model.ValueString(),
		Email:          a.Email.ValueString(),
		Description:    a.Description.ValueString(),
		AiInstructions: a.Instructions.ValueString(),

		AllowedUsers:         users,
		AllowedGroups:        groups,
		Runners:              runners,
		Secrets:              secrets,
		Integrations:         integrations,
		EnvironmentVariables: envVariables,
		Links:                stringList(a.Links.ValueString()),
		Starters:             starters,
	}, err
}

func fromAgent(a *agent, cs *state) (*entities.AgentModel, error) {
	const (
		sep = ","
	)

	by := ""
	at := ""
	email := ""
	var err error
	var intList []string
	var userList []string
	var groupList []string
	var runnerList []string
	var secretList []string
	var starters string

	if len(a.Starters) > 0 {
		jsonBytes, err := json.Marshal(a.Starters)
		if err != nil {
			starters = ""
		}
		starters = string(jsonBytes)
	}

	if a.Metadata != nil {
		at = a.Metadata.CreatedAt
		by = a.Metadata.UserCreated

		for _, u := range cs.users {
			if equal(u.UUID, a.Metadata.UserCreated) {
				email = u.Email
				break
			}
		}
	}

	for _, ar := range a.Runners {
		for _, r := range cs.runners {
			if equal(r.Name, ar) {
				runnerList = append(runnerList, r.Name)
				break
			}
		}
	}
	for _, as := range a.Secrets {
		for _, s := range cs.secrets {
			if equal(s.Name, as) {
				secretList = append(secretList, s.Name)
				break
			}
		}
	}
	for _, ai := range a.Integrations {
		for _, i := range cs.integrations {
			if equal(i.Name, ai) {
				intList = append(intList, i.Name)
				break
			}
		}
	}
	for _, au := range a.AllowedUsers {
		for _, u := range cs.users {
			if equal(u.UUID, au) {
				userList = append(userList, u.Email)
				break
			}
		}
	}
	for _, ag := range a.AllowedGroups {
		for _, g := range cs.groups {
			if equal(g.UUID, ag) {
				groupList = append(groupList, g.Name)
				break
			}
		}
	}

	result := &entities.AgentModel{
		CreatedAt:    types.StringValue(at),
		CreatedBy:    types.StringValue(by),
		Email:        types.StringValue(email),
		Id:           types.StringValue(a.Uuid),
		Name:         types.StringValue(a.Name),
		Image:        types.StringValue(a.Image),
		Model:        types.StringValue(a.LlmModel),
		Description:  types.StringValue(a.Description),
		Instructions: types.StringValue(a.AiInstructions),
		Runners:      types.StringValue(strings.Join(runnerList, sep)),

		Links:        types.StringValue(strings.Join(a.Links, sep)),
		Integrations: types.StringValue(strings.Join(intList, sep)),
		Users:        types.StringValue(strings.Join(userList, sep)),
		Starters:     types.StringValue(starters),
		Groups:       types.StringValue(strings.Join(groupList, sep)),
		Secrets:      types.StringValue(strings.Join(secretList, sep)),
		Variables:    types.StringValue(""),
	}

	if len(a.EnvironmentVariables) >= 1 {
		b, err := json.Marshal(a.EnvironmentVariables)
		if err != nil {
			return nil, err
		}

		result.Variables = types.StringValue(string(b))
	}

	return result, err
}

func (c *Client) ReadAgent(_ context.Context, entity *entities.AgentModel) error {
	if entity != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := entity.Id.ValueString()
		name := entity.Name.ValueString()

		for _, a := range cs.agents {
			if equal(a.Uuid, id) || equal(a.Name, name) {
				entity, err = fromAgent(a, cs)
				break
			}
		}

		return err
	}

	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) DeleteAgent(ctx context.Context, entity *entities.AgentModel) error {
	if entity != nil {
		id := entity.Id.ValueString()
		path := format("/api/v1/agents/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		return err
	}

	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) UpdateAgent(ctx context.Context, entity *entities.AgentModel) error {
	if entity != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := entity.Id.ValueString()
		entity.Email = types.StringValue("")
		uri := c.uri(format("/api/v1/agents/%s", id))

		data, err := toAgent(entity, cs)
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

		entity, err = fromAgent(r, cs)

		return err
	}
	return fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

func (c *Client) CreateAgent(ctx context.Context, entity *entities.AgentModel) (*entities.AgentModel, error) {
	if entity != nil {
		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		data, err := toAgent(entity, cs)
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

	return nil, fmt.Errorf("param entity (*entities.AgentModel) is nil")
}

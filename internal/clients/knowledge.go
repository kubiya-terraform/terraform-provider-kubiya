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

type knowledge struct {
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Labels                []string `json:"labels"`
	Content               string   `json:"content"`
	Groups                []string `json:"groups"`
	Owner                 string   `json:"owner"`
	Type                  string   `json:"type"`
	Source                string   `json:"source"`
	SupportedAgents       []string `json:"supported_agents"`
	SupportedAgentsGroups []string `json:"supported_agents_groups"`
	Id                    string   `json:"uuid"`
}

func toKnowledge(a *entities.KnowledgeModel, cs *state) (*knowledge, error) {
	var err error

	result := &knowledge{
		Source:                "terraform",
		SupportedAgents:       make([]string, 0),
		SupportedAgentsGroups: make([]string, 0),
		Id:                    a.Id.ValueString(),
		Name:                  a.Name.ValueString(),
		Type:                  a.Type.ValueString(),
		Content:               a.Content.ValueString(),
		Description:           a.Description.ValueString(),
	}

	if !a.Labels.IsNull() && !a.Labels.IsUnknown() {
		result.Labels = make([]string, 0)
		for _, v := range a.Labels.Elements() {
			if !v.IsNull() && !v.IsUnknown() {
				str := v.String()
				result.Labels = append(result.Labels, strings.ReplaceAll(str, "\"", ""))
			}
		}
	}

	if !a.Groups.IsNull() && !a.Groups.IsUnknown() {
		result.Groups = make([]string, 0)
		for _, v := range a.Groups.Elements() {
			if !v.IsNull() && !v.IsUnknown() {
				found := false
				str := v.String()
				item := strings.ReplaceAll(str, "\"", "")
				for _, i := range cs.groupList {
					byId := equal(i.UUID, item)
					byName := equal(i.Name, item)
					if found = byId || byName; found {
						result.Groups = append(result.Groups, i.UUID)
						break
					}
				}
				if !found {
					err = errors.Join(err, fmt.Errorf("group \"%s\" don't exist", v))
				}
			}
		}
	}

	if !a.SupportedAgents.IsNull() && !a.SupportedAgents.IsUnknown() {
		result.SupportedAgents = make([]string, 0)
		for _, v := range a.SupportedAgents.Elements() {
			if !v.IsNull() && !v.IsUnknown() {
				found := false
				str := v.String()
				item := strings.ReplaceAll(str, "\"", "")
				for _, i := range cs.agentList {
					byId := equal(i.Uuid, item)
					byName := equal(i.Name, item)
					if found = byId || byName; found {
						result.SupportedAgents = append(result.SupportedAgents, i.Uuid)
						break
					}
				}
				if !found {
					err = errors.Join(err, fmt.Errorf("agent \"%s\" don't exist", item))
				}
			}
		}
	}

	return result, err
}

func fromKnowledge(a *knowledge, cs *state) (*entities.KnowledgeModel, error) {
	var err error
	result := &entities.KnowledgeModel{
		Id:          types.StringValue(a.Id),
		Name:        types.StringValue(a.Name),
		Type:        types.StringValue(a.Type),
		Content:     types.StringValue(a.Content),
		Labels:      toListStringType(a.Labels, err),
		Description: types.StringValue(a.Description),
	}

	if len(a.Groups) >= 1 {
		list := make([]string, 0)
		for _, t := range a.Groups {
			for _, g := range cs.groupList {
				if equal(g.UUID, t) {
					list = append(list, g.Name)
					break
				}
			}
		}
		result.Groups = toListStringType(list, err)
	}

	if len(a.SupportedAgents) >= 1 {
		list := make([]string, 0)
		for _, t := range a.SupportedAgents {
			for _, agentItem := range cs.agentList {
				if equal(agentItem.Uuid, t) {
					list = append(list, agentItem.Name)
					break
				}
			}
		}
		result.SupportedAgents = toListStringType(list, err)
	}

	return result, err
}

func (c *Client) ReadKnowledge(_ context.Context, e *entities.KnowledgeModel) error {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := e.Id
		name := e.Name

		for _, a := range cs.knowledgeList {
			if equal(a.Id, id.ValueString()) ||
				equal(a.Name, name.ValueString()) {
				e, err = fromKnowledge(a, cs)
				break
			}
		}

		return err
	}

	return fmt.Errorf("param entity (*entities.KnowledgeModel) is nil")
}

func (c *Client) DeleteKnowledge(ctx context.Context, e *entities.KnowledgeModel) error {
	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/knowledge/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		return err
	}

	return fmt.Errorf("param entity (*entities.KnowledgeModel) is nil")
}

func (c *Client) UpdateKnowledge(ctx context.Context, e *entities.KnowledgeModel) error {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return err
		}

		id := e.Id.ValueString()
		uri := c.uri(format("/api/v1/knowledge/%s", id))

		data, err := toKnowledge(e, cs)
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

		var r *knowledge
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return err
		}

		e, err = fromKnowledge(r, cs)
		return err
	}
	return fmt.Errorf("param entity (*entities.KnowledgeModel) is nil")
}

func (c *Client) CreateKnowledge(ctx context.Context, e *entities.KnowledgeModel) (*entities.KnowledgeModel, error) {
	if e != nil {
		cs, err := c.state()
		if err != nil {
			return nil, err
		}

		data, err := toKnowledge(e, cs)
		if err != nil {
			return nil, err
		}

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/knowledge")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		var r *knowledge
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		return fromKnowledge(r, cs)
	}

	return e, fmt.Errorf("param entity (*entities.KnowledgeModel) is nil")
}

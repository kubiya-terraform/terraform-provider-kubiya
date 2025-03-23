package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

// runner is used internally by the client
type runner struct {
	Name                string  `json:"name"`
	Subject             string  `json:"subject"`
	Namespace           string  `json:"namespace"`
	RunnerType          string  `json:"runner_type"`
	UserKeyId           string  `json:"user_key_id"`
	Version             int     `json:"version"`
	Description         string  `json:"description"`
	AuthenticationType  string  `json:"authentication_type"`
	ManagedBy           string  `json:"managed_by"`
	TaskId              string  `json:"task_id"`
	WssUrl              string  `json:"wss_url"`
	KubernetesNamespace string  `json:"kubernetes_namespace"`
	GatewayUrl          *string `json:"gateway_url"`
	GatewayPassword     *string `json:"gateway_password"`
	AgentManagerHealth  struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"agent_manager_health"`
	RunnerHealth struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"runner_health"`
	ToolManagerHealth struct {
		Error   string `json:"error"`
		Health  string `json:"health"`
		Status  string `json:"status"`
		Version string `json:"version"`
	} `json:"tool_manager_health"`
}

func (c *Client) ReadRunner(ctx context.Context, entity *entities.RunnerModel) error {
	if entity == nil {
		return fmt.Errorf("param entity (*entities.RunnerModel) is nil")
	}

	const uri = "/api/v3/runners/%s/describe"
	name := entity.Name.ValueString()

	resp, err := c.read(ctx, c.uri(format(uri, name)))
	if err != nil {
		return err
	}

	var r runner
	if err := json.NewDecoder(resp).Decode(&r); err != nil {
		return err
	}

	entity.Name = types.StringValue(r.Name)
	entity.RunnerType = types.StringValue(r.RunnerType)

	return nil
}

func (c *Client) DeleteRunner(ctx context.Context, entity *entities.RunnerModel) error {
	if entity != nil {
		const (
			uri = "/api/v3/runners/%s"
		)
		name := entity.Name.ValueString()

		_, err := c.delete(ctx, c.uri(format(uri, name)))
		return err
	}

	return fmt.Errorf("param entity (*entities.RunnerModel) is nil")
}

func (c *Client) CreateRunner(ctx context.Context, entity *entities.RunnerModel) (*entities.RunnerModel, error) {
	if entity == nil {
		return nil, fmt.Errorf("param entity (*entities.RunnerModel) is nil")
	}

	const uri = "/api/v3/runners/%s"
	name := entity.Name.ValueString()

	data := struct {
		ManagedBy string `json:"managed_by,omitempty"`
		TaskId    string `json:"task_id,omitempty"`
	}{}
	data.ManagedBy, data.TaskId = managedBy()

	body, err := toJson(data)
	if err != nil {
		return nil, err
	}

	reqUri := c.uri(format(uri, name))
	_, err = c.create(ctx, reqUri, body)
	if err != nil {
		return nil, err
	}

	// Now call describe to get the runner type
	if err := c.ReadRunner(ctx, entity); err != nil {
		return nil, fmt.Errorf("runner created but failed to read details: %v", err)
	}

	return entity, nil
}

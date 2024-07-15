package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
)

type scheduledTask struct {
	Id                string                 `json:"task_id"`
	UUID              string                 `json:"task_uuid"`
	Email             string                 `json:"user_email,omitempty"`
	ChannelId         string                 `json:"channel_id,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Agent             string                 `json:"agent,omitempty"`
	TaskType          string                 `json:"task_type,omitempty"`
	ScheduledTime     time.Time              `json:"scheduled_time,omitempty"`
	Status            string                 `json:"status,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	NextScheduledTime time.Time              `json:"next_schedule_time,omitempty"`
}

type createScheduledTaskRequest struct {
	ChannelId     string    `json:"channel_id"`
	ScheduledTime time.Time `json:"schedule_time"`
	Agent         string    `json:"selected_agent"`
	Description   string    `json:"task_description"`
}

func newScheduledTask(body io.Reader) (*scheduledTask, error) {
	var result *scheduledTask
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func fromScheduledTask(a *scheduledTask) (*entities.ScheduledTaskModel, error) {
	var err error
	result := &entities.ScheduledTaskModel{
		Id:                types.StringValue(a.Id),
		UUID:              types.StringValue(a.UUID),
		Status:            types.StringValue(a.Status),
		TaskType:          types.StringValue(a.TaskType),
		ChannelId:         types.StringValue(a.ChannelId),
		Description:       types.StringValue(a.Description),
		ScheduledTime:     types.StringValue(a.ScheduledTime.Format("2006-01-02T15:04:05")),
		NextScheduledTime: types.StringValue(a.NextScheduledTime.Format("2006-01-02T15:04:05")),
	}

	parameters := map[string]string{}
	for key, val := range a.Parameters {
		if key == "context" {
			if str, ok := val.(string); ok {
				result.Agent = types.StringValue(str)
			}
		}

		if str, ok := val.(string); ok {
			parameters[key] = str
		}
	}
	result.Parameters = toMapType(parameters, err)

	return result, err
}

func createScheduledTask(e *entities.ScheduledTaskModel) (*createScheduledTaskRequest, error) {
	var err error
	result := &createScheduledTaskRequest{
		ScheduledTime: time.Time{},
		Agent:         e.Agent.ValueString(),
		ChannelId:     e.ChannelId.ValueString(),
		Description:   e.Description.ValueString(),
	}

	scheduledTime := e.ScheduledTime.ValueString()
	parseTime := func(t, f string, e error) (time.Time, error) {
		if e != nil {
			return time.Time{}, e
		}

		layout := "2006-01-02T15:04:05"
		ts, parseError := time.Parse(layout, t)
		if parseError != nil {
			return time.Time{}, parseError
		}

		if ts.IsZero() {
			return time.Time{}, eformat("%s: %s is not valid time.", f, t)
		}

		return ts, nil
	}
	result.ScheduledTime, err = parseTime(scheduledTime, "scheduled_time", err)

	return result, err
}

func (c *Client) DeleteScheduledTask(ctx context.Context, e *entities.ScheduledTaskModel) error {
	if e != nil {
		id := e.Id.ValueString()
		path := format("/api/v1/scheduled_tasks/%s", id)

		_, err := c.delete(ctx, c.uri(path))
		return err
	}

	return fmt.Errorf("param entity (*entities.ScheduledTaskModel) is nil")
}

func (c *Client) ReadScheduledTask(ctx context.Context, id string) (*entities.ScheduledTaskModel, error) {
	path := format("/api/v1/scheduled_tasks/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		return nil, err
	}

	r, err := newScheduledTask(resp)
	if err != nil {
		return nil, err
	}

	entity, err := fromScheduledTask(r)
	if err != nil || entity == nil {
		if err != nil {
			return nil, err
		}
		return nil, eformat("ScheduledTask %s not found", id)
	}

	return entity, nil
}

func (c *Client) CreateScheduledTask(ctx context.Context, e *entities.ScheduledTaskModel) (*entities.ScheduledTaskModel, error) {
	if e != nil {
		data, err := createScheduledTask(e)
		if err != nil {
			return nil, err
		}

		body, err := toJson(data)
		if err != nil {
			return nil, err
		}

		uri := c.uri("/api/v1/scheduled_tasks")

		resp, err := c.create(ctx, uri, body)
		if err != nil {
			return nil, err
		}

		var r *scheduledTask
		err = json.NewDecoder(resp).Decode(&r)
		if err != nil {
			return nil, err
		}

		return fromScheduledTask(r)
	}

	return e, fmt.Errorf("param entity (*entities.ScheduledTaskModel) is nil")
}

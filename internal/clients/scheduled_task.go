package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya/internal/entities"
	kubiyasentry "terraform-provider-kubiya/internal/sentry"
)

const (
	empty   = ""
	daily   = "daily"
	hourly  = "hourly"
	weekly  = "weekly"
	monthly = "monthly"
)

var (
	cronOptions = []string{daily, hourly, weekly, monthly}
)

type scheduledTask struct {
	Id                string                 `json:"task_id"`
	UUID              string                 `json:"task_uuid"`
	Email             string                 `json:"user_email,omitempty"`
	ChannelId         string                 `json:"channel_id,omitempty"`
	ChannelName       string                 `json:"channel_name,omitempty"`
	Description       string                 `json:"task_description,omitempty"`
	Agent             string                 `json:"agent,omitempty"`
	TaskType          string                 `json:"task_type,omitempty"`
	ScheduledTime     string                 `json:"scheduled_time,omitempty"`
	Status            string                 `json:"status,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	NextScheduledTime string                 `json:"next_schedule_time,omitempty"`
}

type createScheduledTaskRequest struct {
	Email         string    `json:"user_email"`
	ChannelId     string    `json:"channel_id"`
	CronString    string    `json:"cron_string"`
	ScheduledTime time.Time `json:"schedule_time"`
	Agent         string    `json:"selected_agent"`
	Description   string    `json:"task_description"`
	Org           string    `json:"organization_name"`
}

func toDailyCron(t time.Time) string {
	const layout = "%d %d * * * *"
	return fmt.Sprintf(layout, t.Minute(), t.Hour())
}

func toHourlyCron(t time.Time) string {
	const layout = "%d * * * * *"
	return fmt.Sprintf(layout, t.Minute())
}

func toWeeklyCron(t time.Time) string {
	const layout = "%d %d * * %d *"
	return fmt.Sprintf(layout, t.Minute(), t.Hour(), int(t.Weekday()))
}

func toMonthlyCron(t time.Time) string {
	const layout = "%d %d %d * * *"
	return fmt.Sprintf(layout, t.Minute(), t.Hour(), t.Day())
}

func newScheduledTask(body io.Reader) (*scheduledTask, error) {
	var result scheduledTask
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fromScheduledTask(a *scheduledTask) (*entities.ScheduledTaskModel, error) {
	const (
		layout = "2006-01-02T15:04:05"
	)

	var err error
	result := &entities.ScheduledTaskModel{
		Id:                types.StringValue(a.Id),
		UUID:              types.StringValue(a.UUID),
		Email:             types.StringValue(a.Email),
		Repeat:            types.StringValue(""),
		Status:            types.StringValue(a.Status),
		TaskType:          types.StringValue(a.TaskType),
		ChannelId:         types.StringValue(a.ChannelName),
		ScheduledTime:     types.StringValue(a.ScheduledTime),
		NextScheduledTime: types.StringValue(a.NextScheduledTime),
	}

	cron := ""
	parameters := map[string]string{}
	for key, val := range a.Parameters {
		if key == "repeat" {
			if boolean, ok := val.(bool); ok && boolean {
				if _, ok = a.Parameters["cron_string"]; ok {
					item := a.Parameters["cron_string"]
					if _, ok = a.Parameters["cron_string"].(string); ok {
						cron = item.(string)
						if e := result.ParseCron(item.(string)); e != nil {
							err = errors.Join(err, e)
							continue
						}
					}
				}
			}
		}

		if key == "context" {
			if str, ok := val.(string); ok {
				result.Agent = types.StringValue(str)
			}
		}

		if key == "message_text" {
			if str, ok := val.(string); ok {
				result.Description = types.StringValue(str)
			}
		}

		if str, ok := val.(string); ok {
			parameters[key] = str
		}
	}
	result.Parameters = toMapType(parameters, err)

	if t, _ := time.Parse(layout, a.ScheduledTime); t.IsZero() {
		result.Repeat = types.StringValue(cron)
		result.ScheduledTime = types.StringValue(empty)
	}

	return result, err
}

func createScheduledTask(e *entities.ScheduledTaskModel) (*createScheduledTaskRequest, error) {
	var err error

	result := &createScheduledTaskRequest{
		CronString:  empty,
		Agent:       e.Agent.ValueString(),
		ChannelId:   e.ChannelId.ValueString(),
		Description: e.Description.ValueString(),
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

	if !result.ScheduledTime.IsZero() && err == nil {
		switch e.Repeat.ValueString() {
		case daily:
			result.CronString = toDailyCron(result.ScheduledTime)
		case hourly:
			result.CronString = toHourlyCron(result.ScheduledTime)
		case weekly:
			result.CronString = toWeeklyCron(result.ScheduledTime)
		case monthly:
			result.CronString = toMonthlyCron(result.ScheduledTime)
		}
	}

	if cron := e.Repeat.ValueString(); !slices.Contains(cronOptions, cron) && len(cron) > 0 {
		err = nil
		result.CronString = cron
	}

	return result, err
}

func (c *Client) DeleteScheduledTask(ctx context.Context, e *entities.ScheduledTaskModel) error {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "DeleteScheduledTask")
		span.SetData("client.resource_type", "scheduled_task")
	}

	if e == nil {
		logger.Error("DeleteScheduledTask called with nil entity")
		err := fmt.Errorf("param entity (*entities.ScheduledTaskModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	taskId := e.Id.ValueString()
	taskAgent := e.Agent.ValueString()

	logger.Info("Starting scheduled task deletion",
		"task_id", taskId,
		"agent", taskAgent)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Deleting scheduled task via API", sentry.LevelInfo, map[string]interface{}{
		"method":  "DeleteScheduledTask",
		"task_id": taskId,
		"agent":   taskAgent,
	})

	path := format("/api/v1/scheduled_tasks/%s", taskId)
	_, err := c.delete(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to delete scheduled task",
			"task_id", taskId,
			"agent", taskAgent,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return err
	}

	logger.Info("Scheduled task deleted successfully",
		"task_id", taskId,
		"agent", taskAgent)

	return nil
}

func (c *Client) ReadScheduledTask(ctx context.Context, id string) (*entities.ScheduledTaskModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "ReadScheduledTask")
		span.SetData("client.resource_type", "scheduled_task")
	}

	logger.Debug("Reading scheduled task",
		"task_id", id)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Reading scheduled task via API", sentry.LevelInfo, map[string]interface{}{
		"method":  "ReadScheduledTask",
		"task_id": id,
	})

	path := format("/api/v1/scheduled_tasks/%s", id)

	resp, err := c.read(ctx, c.uri(path))
	if err != nil {
		logger.Error("Failed to read scheduled task",
			"task_id", id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	r, err := newScheduledTask(resp)
	if err != nil {
		logger.Error("Failed to decode scheduled task response",
			"task_id", id,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	entity, err := fromScheduledTask(r)
	if err != nil || entity == nil {
		if err != nil {
			logger.Error("Failed to convert scheduled task",
				"task_id", id,
				"error", err.Error())
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		logger.Error("Scheduled task not found",
			"task_id", id)
		err = eformat("ScheduledTask %s not found", id)
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	logger.Debug("Scheduled task read successfully",
		"task_id", id,
		"agent", entity.Agent.ValueString(),
		"status", entity.Status.ValueString())

	return entity, nil
}

func (c *Client) CreateScheduledTask(ctx context.Context, e *entities.ScheduledTaskModel) (*entities.ScheduledTaskModel, error) {
	// Get logger from context
	logger := kubiyasentry.LoggerFromContext(ctx)
	if logger == nil {
		logger = kubiyasentry.GetLogger()
		ctx = kubiyasentry.ContextWithLogger(ctx, logger)
	}

	// Continue tracing from provider level
	span := kubiyasentry.SpanFromContext(ctx)
	if span != nil {
		span.SetData("client.method", "CreateScheduledTask")
		span.SetData("client.resource_type", "scheduled_task")
	}

	if e == nil {
		logger.Error("CreateScheduledTask called with nil entity")
		err := fmt.Errorf("param entity (*entities.ScheduledTaskModel) is nil")
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	taskAgent := e.Agent.ValueString()
	taskDescription := e.Description.ValueString()
	taskChannelId := e.ChannelId.ValueString()

	logger.Info("Starting scheduled task creation",
		"agent", taskAgent,
		"description", taskDescription,
		"channel_id", taskChannelId)

	// Add breadcrumb for client operation
	kubiyasentry.AddBreadcrumb("client", "Creating scheduled task via API", sentry.LevelInfo, map[string]interface{}{
		"method":      "CreateScheduledTask",
		"uri":         "/api/v1/scheduled_tasks",
		"agent":       taskAgent,
		"channel_id":  taskChannelId,
		"description": taskDescription,
	})

	data, err := createScheduledTask(e)
	if err != nil {
		logger.Error("Failed to prepare scheduled task data",
			"agent", taskAgent,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	body, err := toJson(data)
	if err != nil {
		logger.Error("Failed to marshal scheduled task request",
			"agent", taskAgent,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	uri := c.uri("/api/v1/scheduled_tasks")

	resp, err := c.create(ctx, uri, body)
	if err != nil {
		logger.Error("Failed to create scheduled task via API",
			"agent", taskAgent,
			"uri", uri,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	tmp := map[string]string{}
	err = json.NewDecoder(resp).Decode(&tmp)
	if err != nil {
		logger.Error("Failed to decode scheduled task creation response",
			"agent", taskAgent,
			"error", err.Error())
		kubiyasentry.RecordError(ctx, err)
		return nil, err
	}

	if id, ok := tmp["task_id"]; ok {
		logger.Info("Scheduled task created, reading details",
			"task_id", id,
			"agent", taskAgent)

		var entity *entities.ScheduledTaskModel
		if entity, err = c.ReadScheduledTask(ctx, id); err != nil {
			logger.Error("Failed to read scheduled task after creation",
				"task_id", id,
				"agent", taskAgent,
				"error", err.Error())
			kubiyasentry.RecordError(ctx, err)
			return nil, err
		}

		if cron := e.Repeat.ValueString(); !slices.Contains(cronOptions, cron) && len(cron) > 0 {
			entity.Repeat = types.StringValue(cron)
		}

		logger.Info("Scheduled task created successfully",
			"task_id", id,
			"agent", taskAgent,
			"status", entity.Status.ValueString())

		return entity, nil
	}

	logger.Error("Failed to extract task_id from response",
		"agent", taskAgent,
		"response", tmp)
	err = eformat("failed to createWithQueryParams scheduled task")
	kubiyasentry.RecordError(ctx, err)
	return nil, err
}

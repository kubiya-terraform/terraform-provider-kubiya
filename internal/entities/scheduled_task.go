package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ScheduledTaskModel struct {
	Id        types.String `tfsdk:"id"`
	UUID      types.String `tfsdk:"uuid"`
	Email     types.String `tfsdk:"email"`
	Repeat    types.String `tfsdk:"repeat"` // no_repeat, daily, weekly, monthly
	ChannelId types.String `tfsdk:"channel_id"`

	Agent             types.String `tfsdk:"agent"`
	Status            types.String `tfsdk:"status"`
	TaskType          types.String `tfsdk:"task_type"`
	Parameters        types.Map    `tfsdk:"parameters"`
	Description       types.String `tfsdk:"description"`
	ScheduledTime     types.String `tfsdk:"scheduled_time"`
	NextScheduledTime types.String `tfsdk:"next_scheduled_time"`
}

func ScheduledTaskSchema() schema.Schema {
	const (
		empty   = ""
		daily   = "daily"
		hourly  = "hourly"
		weekly  = "weekly"
		monthly = "monthly"
		repeat  = "repeat"
	)

	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Computed
			"id": schema.StringAttribute{
				Computed: true,
			},
			"uuid": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},

			// Required
			"agent": schema.StringAttribute{
				Required: true,
			},
			"channel_id": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},

			// Optional
			"repeat": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(empty),
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"parameters": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"task_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"scheduled_time": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(empty),
			},
			"next_scheduled_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (s *ScheduledTaskModel) ParseCron(cronExpr string) error {
	const (
		daily   = "daily"
		hourly  = "hourly"
		weekly  = "weekly"
		monthly = "monthly"
	)

	if len(cronExpr) <= 0 {
		return nil
	}

	err := validCron(cronExpr)
	if err != nil {
		return err
	}

	switch {
	case isDaily(cronExpr):
		s.Repeat = types.StringValue(daily)
	case isHourly(cronExpr):
		s.Repeat = types.StringValue(hourly)
	case isWeekly(cronExpr):
		s.Repeat = types.StringValue(weekly)
	case isMonthly(cronExpr):
		s.Repeat = types.StringValue(monthly)
	default:
		s.Repeat = types.StringValue(cronExpr)
	}

	return nil
}

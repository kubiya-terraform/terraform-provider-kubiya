package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ScheduledTaskModel struct {
	Id        types.String `tfsdk:"id"`
	UUID      types.String `tfsdk:"uuid"`
	Email     types.String `tfsdk:"email"`
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
			"scheduled_time": schema.StringAttribute{
				Required: true,
			},

			// Optional
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
			"next_scheduled_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

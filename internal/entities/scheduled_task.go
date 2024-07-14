package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ScheduledTaskModel struct {
	Id        types.String `tfsdk:"id"`
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
				Computed:            true,
				Description:         "The ID of the agent",
				MarkdownDescription: "The unique identifier of the agent",
			},

			// Required
			"email": schema.StringAttribute{
				Computed:            true,
				Description:         "The owner of the agent",
				MarkdownDescription: "The user who created the agent",
			},
			"channel_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation time of the agent",
				MarkdownDescription: "The timestamp when the agent was created",
			},

			// Optional
			"agent": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
			},
			"parameters": schema.MapAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"task_type": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"scheduled_time": schema.StringAttribute{
				Optional: true,
			},
			"next_scheduled_time": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

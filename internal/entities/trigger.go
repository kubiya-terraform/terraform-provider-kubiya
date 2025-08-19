package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TriggerModel represents the Terraform resource model for a trigger.
type TriggerModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Runner     types.String `tfsdk:"runner"`
	Workflow   types.String `tfsdk:"workflow"`
	Url        types.String `tfsdk:"url"`
	Status     types.String `tfsdk:"status"`
	WorkflowId types.String `tfsdk:"workflow_id"`
}

// TriggerSchema defines the schema for the trigger resource.
func TriggerSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a Kubiya workflow trigger with webhook capabilities",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the trigger",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the trigger",
			},
			"runner": schema.StringAttribute{
				Required:    true,
				Description: "Runner to use for executing the workflow (e.g., 'kubiya-hosted', 'core-testing-1')",
			},
			"workflow": schema.StringAttribute{
				Required:    true,
				Description: "JSON-encoded workflow definition containing name, version, and steps",
				PlanModifiers: []planmodifier.String{
					jsonNormalizationModifier(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The webhook URL for triggering the workflow",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the workflow (e.g., 'draft', 'published')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workflow_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the created workflow",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

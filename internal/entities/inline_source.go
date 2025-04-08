package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ArgModel struct {
	Name        types.String     `tfsdk:"name"`
	Type        types.String     `tfsdk:"type"`
	Default     types.String     `tfsdk:"default"`
	Options     types.List       `tfsdk:"options"`
	Required    types.Bool       `tfsdk:"required"`
	Description types.String     `tfsdk:"description"`
	OptionsFrom OptionsFormModel `tfsdk:"options_from"`
}

type OptionsFormModel struct {
	Image  types.String `tfsdk:"image"`
	Script types.String `tfsdk:"script"`
}
type FileSpecModel struct {
	Source      types.String `tfsdk:"source"`
	Content     types.String `tfsdk:"content"`
	Destination types.String `tfsdk:"destination"`
}

type InlineTool struct {
	Icon        types.String `tfsdk:"icon"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Image       types.String `tfsdk:"image"`
	Content     types.String `tfsdk:"content"`
	Mermaid     types.String `tfsdk:"mermaid"`
	OnStart     types.String `tfsdk:"on_start"`
	OnBuild     types.String `tfsdk:"on_build"`
	Description types.String `tfsdk:"description"`
	OnComplete  types.String `tfsdk:"on_complete"`

	Env        types.List      `tfsdk:"env"`
	Args       []ArgModel      `tfsdk:"args"`
	Files      []FileSpecModel `tfsdk:"files"`
	Secrets    types.List      `tfsdk:"secrets"`
	Entrypoint types.List      `tfsdk:"entrypoint"`

	Workflow    types.Bool `tfsdk:"workflow"`
	LongRunning types.Bool `tfsdk:"long_running"`
}

type InlineSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Tools  []InlineTool `tfsdk:"tools"`
	Runner types.String `tfsdk:"runner"`
	Config types.String `tfsdk:"dynamic_config"`
}

func InlineSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Computed
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the tool",
				MarkdownDescription: "The unique identifier of the inline source tool",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the inline source",
				MarkdownDescription: "The descriptive type of the inline source",
			},

			// Required
			"tools": schema.ListNestedAttribute{
				Required:    true,
				Description: "A list of tools for inline source",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"workflow": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"long_running": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},

						"icon":        schema.StringAttribute{Optional: true},
						"name":        schema.StringAttribute{Required: true},
						"type":        schema.StringAttribute{Computed: true, Optional: true, Default: defaultString()},
						"image":       schema.StringAttribute{Optional: true},
						"content":     schema.StringAttribute{Optional: true},
						"mermaid":     schema.StringAttribute{Optional: true},
						"on_start":    schema.StringAttribute{Optional: true},
						"on_build":    schema.StringAttribute{Optional: true},
						"description": schema.StringAttribute{Required: true},
						"on_complete": schema.StringAttribute{Optional: true},

						"env":        schema.ListAttribute{Optional: true, ElementType: types.StringType},
						"secrets":    schema.ListAttribute{Optional: true, ElementType: types.StringType},
						"entrypoint": schema.ListAttribute{Optional: true, ElementType: types.StringType},

						"args": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"name":        schema.StringAttribute{Required: true},
							"type":        schema.StringAttribute{Computed: true, Optional: true, Default: defaultString()},
							"description": schema.StringAttribute{Required: true},
							"required":    schema.BoolAttribute{Optional: true},
							"default":     schema.StringAttribute{Optional: true},
							"options":     schema.ListAttribute{Optional: true, ElementType: types.StringType},
							"options_from": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"image":  schema.StringAttribute{Required: true},
									"script": schema.StringAttribute{Required: true},
								},
							},
						}}},
						"files": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
							"source":      schema.StringAttribute{Optional: true},
							"destination": schema.StringAttribute{Required: true},
							"content":     schema.StringAttribute{Optional: true},
						}}},
					},
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the inline source tool",
				MarkdownDescription: "The descriptive name of the inline source",
			},

			// Optional + Computed
			"runner": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "The runner name",
				MarkdownDescription: "The runner name to add for inline source",
			},
			"dynamic_config": schema.StringAttribute{
				Optional:            true,
				PlanModifiers:       []planmodifier.String{jsonNormalizationModifier()},
				Description:         "The dynamic configuration of the inline source",
				MarkdownDescription: "A map of key-value pairs representing dynamic configuration for the inline source",
			},
		},
	}
}

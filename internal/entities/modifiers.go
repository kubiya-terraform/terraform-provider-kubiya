package entities

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = jsonStringModifier{}

type jsonStringModifier struct{}

func (j jsonStringModifier) Description(_ context.Context) string {
	return "modify json string"
}

func (j jsonStringModifier) MarkdownDescription(_ context.Context) string {
	return "modify json string"
}

func (j jsonStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	tmp := map[string]interface{}{}
	planValue := req.PlanValue.ValueString()
	err := json.Unmarshal([]byte(planValue), &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"json unmarshal error",
			err.Error(),
		)
		return
	}

	planResult, err := json.Marshal(tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"json marshal error",
			err.Error(),
		)
		return
	}

	resp.PlanValue = types.StringValue(string(planResult))
}

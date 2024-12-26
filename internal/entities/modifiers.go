package entities

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ planmodifier.String = &jsonStringModifier{}
)

type jsonStringModifier struct{}

func jsonNormalizationModifier() planmodifier.String {
	modifier := &jsonStringModifier{}
	return modifier
}

func (j *jsonStringModifier) Description(_ context.Context) string {
	return "Normalizes JSON string format for consistent comparison"
}

func (j *jsonStringModifier) MarkdownDescription(_ context.Context) string {
	return "Normalizes JSON string format for consistent comparison"
}

func (j *jsonStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	// Normalize both config and plan values
	tmp := map[string]interface{}{}
	configValue := req.ConfigValue.ValueString()
	err := json.Unmarshal([]byte(configValue), &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing JSON",
			"Failed to parse JSON string: "+err.Error(),
		)
		return
	}

	normalizedResult, err := json.Marshal(tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error normalizing JSON",
			"Failed to normalize JSON: "+err.Error(),
		)
		return
	}

	// Set the normalized config value as the plan value
	resp.PlanValue = types.StringValue(string(normalizedResult))
}

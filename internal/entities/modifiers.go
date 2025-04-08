package entities

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

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
	configValue := req.ConfigValue.ValueString()
	normalizedResult, err := normalizeJSON(configValue)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error normalizing JSON",
			"Failed to normalize JSON: "+err.Error(),
		)
		return
	}

	// Set the normalized config value as the plan value
	resp.PlanValue = types.StringValue(normalizedResult)
}

func normalizeJSON(input string) (string, error) {
	var data interface{}

	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}

	// Marshal with sorted keys and consistent output
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "") // no extra whitespace
	if err := enc.Encode(data); err != nil {
		return "", err
	}

	// Remove trailing newline added by Encoder
	return strings.TrimSpace(buf.String()), nil
}

package clients

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

// Function to process types.Map and convert it to map[string]interface{}
func processDynamicConfig(dynamicConfig types.Map) (map[string]any, error) {
	// Ensure the map is not null or unknown
	if dynamicConfig.IsNull() || dynamicConfig.IsUnknown() {
		return nil, nil
	}

	// Declare the output map
	result := make(map[string]any)

	// Retrieve the elements of the map
	for key, value := range dynamicConfig.Elements() {
		var goValue any

		// Check the type of each element and convert it accordingly
		switch v := value.(type) {
		case types.String:
			goValue = v.ValueString()
		case types.Int64:
			goValue = v.ValueInt64()
		case types.Float64:
			goValue = v.ValueFloat64()
		case types.Bool:
			goValue = v.ValueBool()
		default:
			goValue = fmt.Sprintf("Unsupported type: %T", v)
		}

		result[key] = goValue
	}

	return result, nil
}

// Function to convert map[string]interface{} to types.Map
func convertToTypesMap(input map[string]any) (types.Map, error) {
	// Create a map of attr.Value
	attributeValues := make(map[string]attr.Value)

	// Iterate through the input map and convert each value to attr.Value
	for key, value := range input {
		var attrValue attr.Value

		switch v := value.(type) {
		case string:
			attrValue = types.StringValue(v)
		case int:
			attrValue = types.Int64Value(int64(v))
		case int64:
			attrValue = types.Int64Value(v)
		case float64:
			attrValue = types.Float64Value(v)
		case bool:
			attrValue = types.BoolValue(v)
		case []byte:
			// Encode byte arrays as base64 strings
			attrValue = types.StringValue(base64.StdEncoding.EncodeToString(v))
		default:
			return types.Map{}, fmt.Errorf("unsupported Type %s", fmt.Sprintf("Unsupported type for key '%s': %T", key, v))
		}

		attributeValues[key] = attrValue
	}

	// Create a types.MapValue
	result, diags := types.MapValue(types.DynamicType, attributeValues)
	if diags.HasError() {
		err := ConvertDiagnosticsToError(diags)
		return types.Map{}, err
	}

	return result, nil
}

// ConvertDiagnosticsToError converts diag.Diagnostics to a single error
func ConvertDiagnosticsToError(diags diag.Diagnostics) error {
	if diags == nil || len(diags) == 0 {
		return nil
	}

	// Combine all diagnostic messages into a single string
	var errorMessages []string
	for _, d := range diags {
		errorMessages = append(errorMessages, d.Detail())
	}

	return errors.New(strings.Join(errorMessages, "; "))
}

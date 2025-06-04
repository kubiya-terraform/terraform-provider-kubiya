package vendors

import (
	"context"
	"fmt"
	"math/big"

	"terraform-provider-kubiya/internal/entities"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BaseExternalKnowledge contains common fields for all vendor responses
type BaseExternalKnowledge struct {
	UUID            string `json:"uuid"`
	Org             string `json:"org"`
	StartDate       string `json:"start_date"`
	IntegrationType string `json:"integration_type"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ConvertToTerraformValue converts Go values to Terraform attr.Value
func ConvertToTerraformValue(v interface{}) attr.Value {
	switch val := v.(type) {
	case string:
		return types.DynamicValue(types.StringValue(val))
	case []interface{}:
		// Handle lists
		elements := make([]attr.Value, len(val))
		for i, item := range val {
			if str, ok := item.(string); ok {
				elements[i] = types.StringValue(str)
			} else {
				elements[i] = types.StringValue(fmt.Sprintf("%v", item))
			}
		}
		// Create tuple type for the elements
		elementTypes := make([]attr.Type, len(elements))
		for i := range elements {
			elementTypes[i] = types.StringType
		}
		tupleVal, _ := types.TupleValue(elementTypes, elements)
		return types.DynamicValue(tupleVal)
	case []string:
		// Handle string slices directly - return as tuple to match HCL input
		elements := make([]attr.Value, len(val))
		elementTypes := make([]attr.Type, len(val))
		for i, str := range val {
			elements[i] = types.StringValue(str)
			elementTypes[i] = types.StringType
		}
		tupleVal, _ := types.TupleValue(elementTypes, elements)
		return types.DynamicValue(tupleVal)
	case bool:
		return types.DynamicValue(types.BoolValue(val))
	case float64:
		return types.DynamicValue(types.NumberValue(big.NewFloat(val)))
	case int64:
		return types.DynamicValue(types.NumberValue(big.NewFloat(float64(val))))
	case int:
		return types.DynamicValue(types.NumberValue(big.NewFloat(float64(val))))
	default:
		// Fallback to string representation
		return types.DynamicValue(types.StringValue(fmt.Sprintf("%v", val)))
	}
}

// CreateExternalKnowledgeModel creates a model from base fields and config
func CreateExternalKnowledgeModel(base BaseExternalKnowledge, vendor string, configElements map[string]attr.Value) *entities.ExternalKnowledgeModel {
	// Create an object type based on the config elements
	attrTypes := make(map[string]attr.Type)
	for key, val := range configElements {
		attrTypes[key] = val.Type(context.Background())
	}

	configObj, _ := types.ObjectValue(attrTypes, configElements)

	return &entities.ExternalKnowledgeModel{
		Id:              types.StringValue(base.UUID),
		Vendor:          types.StringValue(vendor),
		Config:          types.DynamicValue(configObj),
		Org:             types.StringValue(base.Org),
		StartDate:       types.StringValue(base.StartDate),
		IntegrationType: types.StringValue(base.IntegrationType),
		CreatedAt:       types.StringValue(base.CreatedAt),
		UpdatedAt:       types.StringValue(base.UpdatedAt),
	}
}

// ExtractDynamicValue extracts the underlying value from a dynamic Terraform value
func ExtractDynamicValue(dynVal types.Dynamic) interface{} {
	if dynVal.IsNull() || dynVal.IsUnknown() {
		return nil
	}

	underlyingVal := dynVal.UnderlyingValue()

	switch val := underlyingVal.(type) {
	case types.String:
		if !val.IsNull() && !val.IsUnknown() {
			return val.ValueString()
		}
	case types.List:
		if !val.IsNull() && !val.IsUnknown() {
			var stringList []string
			for _, elem := range val.Elements() {
				if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					stringList = append(stringList, strVal.ValueString())
				}
			}
			return stringList
		}
	case types.Tuple:
		if !val.IsNull() && !val.IsUnknown() {
			var stringList []string
			for _, elem := range val.Elements() {
				if strVal, ok := elem.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
					stringList = append(stringList, strVal.ValueString())
				}
			}
			return stringList
		}
	case types.Bool:
		if !val.IsNull() && !val.IsUnknown() {
			return val.ValueBool()
		}
	case types.Number:
		if !val.IsNull() && !val.IsUnknown() {
			if floatVal, _ := val.ValueBigFloat().Float64(); floatVal == float64(int64(floatVal)) {
				return int64(floatVal)
			} else {
				return floatVal
			}
		}
	}

	return nil
}

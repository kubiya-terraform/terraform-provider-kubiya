package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func cleanString(str string) string {
	remove := "\""
	return strings.TrimSuffix(strings.TrimPrefix(str, remove), remove)
}

func splitToStringList(str string) []string {
	prefix := "["
	suffix := "]"
	str = strings.TrimPrefix(str, prefix)
	return strings.Split(strings.TrimSuffix(str, suffix), ",")
}

func toStringSlice(item types.List) []string {
	var result []string
	for _, str := range splitToStringList(item.String()) {
		result = append(result, cleanString(str))
	}
	return result
}

func convertTypesMapToStringMap(input types.Map) map[string]string {
	result := make(map[string]string)
	if input.IsNull() || input.IsUnknown() {
		return result
	}

	for key, val := range input.Elements() {
		str, ok := val.(types.String)
		if ok && !str.IsNull() && !str.IsUnknown() {
			result[key] = str.ValueString()
		}
	}

	return result
}

func toListType(diags *diag.Diagnostics, items ...string) types.List {
	var list []attr.Value
	for _, item := range items {
		list = append(list, types.StringValue(cleanString(item)))
	}

	result, d := types.ListValue(types.StringType, list)
	if d.HasError() {
		diags.Append(d.Errors()...)
	}

	return result
}

func convertStringMapToMapType(diags *diag.Diagnostics, input map[string]string) types.Map {
	elements := make(map[string]attr.Value)
	for key, val := range input {
		elements[key] = types.StringValue(val)
	}

	result, d := types.MapValue(types.StringType, elements)
	if d.HasError() {
		diags.Append(d...)
	}

	return result
}

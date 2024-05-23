package entities

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsOneOf struct {
	field string
	list  []string
}

func onOfValidator(f string, l []string) stringIsOneOf {
	return stringIsOneOf{
		field: f,
		list:  append(make([]string, 0), l...),
	}
}

func (v stringIsOneOf) Description(_ context.Context) string {
	return fmt.Sprintf("value must be on of: %s", strings.Join(v.list, ","))
}

func (v stringIsOneOf) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be on of: %s", strings.Join(v.list, ","))
}

func (v stringIsOneOf) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if found := slices.Contains(v.list, value); !found {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("Invalid Value in `%s` field", v.field),
			fmt.Sprintf("invalid value allowd. options: [%s]", strings.Join(v.list, ", ")),
		)
		return
	}
}

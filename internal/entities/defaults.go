package entities

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

func defaultString(i ...string) defaults.String {
	val := ""
	if len(i) > 0 {
		val = i[0]
	}
	return stringdefault.StaticString(val)
}

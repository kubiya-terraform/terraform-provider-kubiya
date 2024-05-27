package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func closeBody(b io.ReadCloser) {
	_ = b.Close()
}

func toPathYaml(pre, suf string) string {
	slash := "/"
	layout := "%s/%s.yaml"
	pre = strings.TrimSuffix(pre, slash)
	suf = strings.TrimPrefix(suf, slash)

	return fmt.Sprintf(layout, pre, suf)
}

func toJson(item interface{}) (io.Reader, error) {
	if item != nil {
		body, err := toJsonBytes(item)
		if err != nil || len(body) <= 0 {
			if err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("item is nil")
		}
		return bytes.NewReader(body), nil
	}

	return nil, fmt.Errorf("item is nil")
}

func toJsonBytes(item interface{}) ([]byte, error) {
	if item != nil {
		body, err := json.Marshal(item)
		if err != nil || len(body) <= 0 {
			if err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("item is nil")
		}
		return body, nil
	}

	return nil, fmt.Errorf("item is nil")
}

func eformat(l string, i ...any) error {
	return fmt.Errorf(l, i...)
}

func format(l string, i ...any) string {
	return fmt.Sprintf(l, i...)
}

func equal(str, term string) bool {
	return strings.EqualFold(str, term)
}

func toListStringType(items []string, err error) types.List {
	const (
		errorLayout = "[%s] %s. Error: %s"
	)

	var elements []attr.Value
	for _, item := range items {
		elements = append(elements, types.StringValue(item))
	}

	result, diags := types.ListValue(types.StringType, elements)
	if diags.HasError() {
		for _, d := range diags {
			detail := d.Detail()
			summary := d.Summary()
			severity := d.Severity()
			err = errors.Join(err, eformat(errorLayout, severity, summary, detail))
		}
	}

	return result
}

func toMapType(items map[string]string, err error) types.Map {
	const (
		errorLayout = "[%s] %s. Error: %s"
	)

	elements := make(map[string]attr.Value)
	for key, value := range items {
		if len(key) >= 1 {
			elements[key] = types.StringValue(value)
		}
	}

	result, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		for _, d := range diags {
			detail := d.Detail()
			summary := d.Summary()
			severity := d.Severity()
			err = errors.Join(err, eformat(errorLayout, severity, summary, detail))
		}
	}

	return result
}

func fromList(items types.List) []string {
	var result []string
	for _, item := range items.Elements() {
		str := item.String()
		result = append(result, strings.ReplaceAll(str, "\"", ""))
	}

	return result
}

package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const errorLayout = "[%s] %s. Error: %s"

func closeBody(b io.ReadCloser) {
	_ = b.Close()
}

func fromJson(r io.Reader, item any) error {
	if r == nil {
		return fmt.Errorf("response body is nil")
	}

	if item == nil {
		return fmt.Errorf("item is nil")
	}

	if err := json.NewDecoder(r).Decode(item); err != nil {
		return err
	}

	return nil
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

func toListStringType(items []string, err error) types.List {
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

func getTaskId() string {
	const (
		empty = ""
		env   = "TASK_UUID"
	)
	if value := os.Getenv(env); value != "" {
		return value
	}

	return empty
}

func getManagedBy() string {
	const (
		byTerraform = "terraform"
		env         = "MANAGED_BY"
	)
	if value := os.Getenv(env); value != "" {
		return value
	}

	return byTerraform
}

func managedBy() (string, string) {
	const (
		defaultId = ""
		defaultBy = "terraform"

		taskIdEnv    = "TASK_UUID"
		managedByEnv = "MANAGED_BY"
	)

	id := defaultId
	by := defaultBy

	if value := os.Getenv(taskIdEnv); value != "" {
		id = value
	}

	if value := os.Getenv(managedByEnv); value != "" {
		by = value
	}

	return by, id
}

func responseBodyError(r *http.Response) error {
	body, err := responseBody(r)
	if err != nil {
		return err
	}

	return fmt.Errorf(body)
}

func responseBody(r *http.Response) (string, error) {
	if r == nil {
		return "", nil
	}

	defer closeBody(r.Body)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func cleanMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		// Skip nil values
		if v == nil {
			continue
		}

		// Skip empty strings
		if s, ok := v.(string); ok && s == "" {
			continue
		}

		// Handle empty maps - either empty or only containing empty values
		if subMap, ok := v.(map[string]interface{}); ok {
			cleaned := cleanMap(subMap)
			if len(cleaned) > 0 {
				result[k] = cleaned
			}
			continue
		}

		// Handle slices/arrays
		if arr, ok := v.([]interface{}); ok {
			if len(arr) == 0 {
				continue
			}

			// Clean each item in the array if it's a map
			cleanedArr := make([]interface{}, 0, len(arr))
			for _, item := range arr {
				if itemMap, ok := item.(map[string]interface{}); ok {
					cleaned := cleanMap(itemMap)
					if len(cleaned) > 0 {
						cleanedArr = append(cleanedArr, cleaned)
					}
				} else {
					// Keep non-map items as is
					cleanedArr = append(cleanedArr, item)
				}
			}

			if len(cleanedArr) > 0 {
				result[k] = cleanedArr
			}
			continue
		}

		// Keep non-empty, non-map, non-slice values
		result[k] = v
	}

	return result
}

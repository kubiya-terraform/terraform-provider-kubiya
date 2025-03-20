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

func clean(str, o, n string) string {
	return strings.ReplaceAll(str, o, n)
}

func toPathYaml(pre, suf string) string {
	slash := "/"
	layout := "%s/%s.yaml"
	pre = strings.TrimSuffix(pre, slash)
	suf = strings.TrimPrefix(suf, slash)

	return fmt.Sprintf(layout, pre, suf)
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

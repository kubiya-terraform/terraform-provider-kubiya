package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func structName(i any) string {
	t := reflect.TypeOf(i)
	return t.Name()
}

func closeBody(b io.ReadCloser) {
	_ = b.Close()
}

func stringList(str string) []string {
	const (
		sep   = ","
		pre   = "["
		suf   = "]"
		empty = ""
	)

	str = strings.ReplaceAll(str, pre, empty)
	return strings.Split(strings.ReplaceAll(str, suf, empty), sep)
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
		body, err := json.Marshal(item)
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

func eformat(l string, i ...any) error {
	return fmt.Errorf(l, i...)
}

func format(l string, i ...any) string {
	return fmt.Sprintf(l, i...)
}

func equal(str, term string) bool {
	return strings.EqualFold(str, term)
}

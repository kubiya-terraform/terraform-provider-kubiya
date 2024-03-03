package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

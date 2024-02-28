package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func closeBody(b io.ReadCloser) {
	_ = b.Close()
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

package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// postJSON 送 JSON、讀回覆（最多 1 MB）。answered 是對方有沒有回應（連不上、逾時是 false）；HTTP 不是 200 也回 error
func postJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, body any) (data []byte, answered bool, err error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, false, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("content-type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, true, err
	}
	if resp.StatusCode != http.StatusOK {
		return data, true, fmt.Errorf("HTTP %d：%s", resp.StatusCode, clip(data))
	}
	return data, true, nil
}

func clip(b []byte) string {
	if len(b) > 300 {
		return string(b[:300]) + "…"
	}
	return string(b)
}

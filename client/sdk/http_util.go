package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const hzmJobAccessToken = "hzm-job-access-token"

type RemotingUtil struct {
	client   *http.Client
	maxRetry int
}

func New(opts ...Option) *RemotingUtil {
	h := &RemotingUtil{
		client:   &http.Client{Timeout: 3 * time.Second},
		maxRetry: 2,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

type Option func(*RemotingUtil)

func WithTimeout(d time.Duration) Option {
	return func(h *RemotingUtil) {
		h.client.Timeout = d
	}
}

func WithMaxRetry(maxRetry int) Option {
	return func(h *RemotingUtil) {
		h.maxRetry = maxRetry
	}
}

func (h *RemotingUtil) PostJSON(ctx context.Context, url, accessToken string, body any) ([]byte, error) {
	method := http.MethodPost
	headers := map[string]string{
		"Content-Type":    "application/json",
		hzmJobAccessToken: accessToken,
	}

	var b io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		b = bytes.NewBuffer(jsonData)
	}
	return h.doRequest(ctx, method, url, b, headers)
}

func (h *RemotingUtil) doRequest(ctx context.Context, method string, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	var lastErr error
	for i := 0; i <= h.maxRetry; i++ {
		req, _ := http.NewRequestWithContext(ctx, method, url, body)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
			continue
		}
		return data, nil
	}
	return nil, lastErr
}

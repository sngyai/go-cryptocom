package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sngyai/go-cryptocom/errors"
)

type Requester struct {
	Client  *http.Client
	BaseURL string
}

func (r Requester) Post(ctx context.Context, body Request, method string, response interface{}) (int, error) {
	return r.doRequest(ctx, http.MethodPost, body, method, response)
}

func (r Requester) Get(ctx context.Context, body Request, method string, response interface{}) (int, error) {
	return r.doRequest(ctx, http.MethodGet, body, method, response)
}

func (r Requester) doRequest(ctx context.Context, httpMethod string, body Request, method string, response interface{}) (int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	version := V1
	if body.Version != "" {
		version = body.Version
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, fmt.Sprintf("%s%s%s", r.BaseURL, version, method), bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := r.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(resBytes, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response body: %s, error: %w", string(resBytes), err)
	}

	return res.StatusCode, nil
}

func (Requester) CheckErrorResponse(statusCode int, responseCode json.Number) error {
	if statusCode >= 400 {
		code, err := responseCode.Int64()
		if err != nil {
			return errors.ResponseError{
				HTTPStatusCode: statusCode,
				Err:            fmt.Errorf("invalid response code: %v", responseCode),
			}
		}
		return errors.NewResponseError(statusCode, code)
	}

	return nil
}

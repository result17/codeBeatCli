package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/result17/codeBeatCli/internal/heartbeat"
	"github.com/result17/codeBeatCli/pkg/log"
)

const (
	CollectHeartbeatRouter = "/users/current/heartbeats/collect"
)

func (c Client) SendHeartbeats(ctx context.Context, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
	logger := log.Extract(ctx)
	logger.Debugf("Sending %d heartbeats(s) to api at %s", len(hs), c.baseURL)

	var results []heartbeat.Result

	collectRouterPath := c.baseURL + CollectHeartbeatRouter

	res, err := c.sendHeartbeats(ctx, collectRouterPath, hs)

	if err != nil {
		return nil, err
	}

	results = append(results, res...)

	return results, nil
}

func (c Client) sendHeartbeats(ctx context.Context, url string, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
	logger := log.Extract(ctx)
	data, err := json.Marshal(hs)

	if err != nil {
		return nil, fmt.Errorf("Failed to json encode heartbeats: %s", err)
	}

	logger.Debugf("heartbeats: %s", string(data))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))

	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("Request to %q timed out", url)
		}
		return nil, fmt.Errorf("Failed making request to %q: %s", url, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed reading response body from %q: %s", url, err)
	}
	switch res.StatusCode {
	case http.StatusAccepted, http.StatusCreated:
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad request at %q", url)
	default:
		return nil, fmt.Errorf("Invalid response status from %q. got: %d, want: %d/%d. body: %q",
			url, res.StatusCode, http.StatusAccepted, http.StatusCreated, string(body))
	}
	result, err := ParseHeartbeatResponses(ctx, body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ParseHeartbeatResponses(ctx context.Context, data []byte) ([]heartbeat.Result, error) {
	var responsesBody struct {
		Responses [][]json.RawMessage `json:"responses"`
	}
	err := json.Unmarshal(data, &responsesBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse json response body: %s. body: %q", err, string(data))
	}

	var results []heartbeat.Result
	for n, r := range responsesBody.Responses {
		result, err := parseHeartbeatResponse(ctx, r)
		if err != nil {
			return nil, fmt.Errorf("Failed parsing result #%d: %s. body: %q", n, err, string(data))
		}
		results = append(results, result)
	}
	return results, nil
}

func parseHeartbeatResponse(ctx context.Context, data []json.RawMessage) (heartbeat.Result, error) {
	var result heartbeat.Result

	type responseBody struct {
		Data *heartbeat.Heartbeat `json:"data"`
	}

	err := json.Unmarshal(data[1], &result.Status)
	if err != nil {
		return heartbeat.Result{}, fmt.Errorf("failed to parse json status: %s", err)
	}

	if result.Status < http.StatusOK || result.Status > 299 {
		resultErrors, err := parseHeartbeatResponseError(ctx, data[0])
		if err != nil {
			return heartbeat.Result{}, fmt.Errorf("failed to parse result errors: %s", err)
		}

		result.Errors = resultErrors

		return heartbeat.Result{
			Errors: result.Errors,
			Status: result.Status,
		}, nil
	}

	err = json.Unmarshal(data[0], &responseBody{Data: &result.Heartbeat})
	if err != nil {
		return heartbeat.Result{}, fmt.Errorf("failed to parse json heartbeat: %s", err)
	}

	return result, nil

}

func parseHeartbeatResponseError(ctx context.Context, data json.RawMessage) ([]string, error) {
	logger := log.Extract(ctx)
	var errs []string
	type responseBodyErr struct {
		Error  *string         `json:"error"`
		Errors *map[string]any `json:"errors"`
	}

	// 1. try "error" key
	var resultError string

	err := json.Unmarshal(data, &responseBodyErr{Error: &resultError})
	if err != nil {
		logger.Debugf("Failed to parse json heartbeat error or 'error' key note found: %s", err)
	}

	if resultError != "" {
		errs = append(errs, resultError)
		return errs, nil
	}

	// 2. try "errors" key
	var resultErrors map[string]any

	err = json.Unmarshal(data, &responseBodyErr{Errors: &resultErrors})
	if err != nil {
		logger.Debugf("Failed to parse json heartbeat error or 'error' key note found: %s", err)
	}
	if resultErrors == nil {
		return nil, errors.New("failed to detect any errors despite invalid response status")
	}

	for field, messages := range resultErrors {
		// skipping parsing dependencies errors as it won't happen because we are
		// filtering in the cli.
		if field == "dependencies" {
			continue
		}

		m := make([]string, len(messages.([]any)))
		for i, v := range messages.([]any) {
			m[i] = fmt.Sprint(v)
		}

		errs = append(errs, fmt.Sprintf(
			"%s: %s",
			field,
			strings.Join(m, " "),
		))
	}
	return errs, nil
}

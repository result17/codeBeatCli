package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/result17/codeBeatCli/internal/heartbeat"
	"github.com/result17/codeBeatCli/pkg/log"
)

const (
	CollectHeartbeatRouter = "/api/heartbeat/list"
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

// sendHeartbeats main logic
func (c Client) sendHeartbeats(ctx context.Context, url string, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
	logger := log.Extract(ctx)
	data, err := json.Marshal(hs)

	if err != nil {
		return nil, fmt.Errorf("Failed to json encode heartbeats: %s", err)
	}

	logger.Debugf("Heartbeats: %s", string(data))

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

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed reading response body from %q: %s", url, err)
	}
	switch res.StatusCode {
	case http.StatusAccepted, http.StatusCreated:
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad request at %q", url)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("Server error at %q", url)
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
	var responsesBody []json.RawMessage
	err := json.Unmarshal(data, &responsesBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse json response body: %s. body: %q", err, string(data))
	}
	var results []heartbeat.Result
	for n, r := range responsesBody {
		result, err := parseHeartbeatResponse(ctx, r)
		if err != nil {
			return nil, fmt.Errorf("Failed parsing result #%d: %s. body: %q", n, err, string(data))
		}
		results = append(results, result)
	}
	return results, nil
}

func parseHeartbeatResponse(ctx context.Context, data json.RawMessage) (heartbeat.Result, error) {
	var result heartbeat.Result

	type responseBody struct {
		Data   *heartbeat.Heartbeat `json:"data"`
		Status int                  `json:"status"`
	}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return heartbeat.Result{}, fmt.Errorf("Failed to parse json status or heartbeat: %s", err)
	}

	if result.Status < http.StatusOK || result.Status > 299 {
		return heartbeat.Result{}, fmt.Errorf("Incorrect status: %d", result.Status)
	}

	return result, nil
}

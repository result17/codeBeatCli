package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/result17/codeBeatCli/internal/summary"
)

const (
	TodayRouter = "/api/duration/today"
)

func (c *Client) Today(ctx context.Context) (*summary.GrandTotal, error) {
	url := c.baseURL + TodayRouter

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	resp, err := c.Do(ctx, req)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to make request to %q: %s", url, err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read body from %q: %s", url, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad request at %q", url)
	default:
		return nil, fmt.Errorf(
			"Invalid response status from %q. want %d, got %d. body: %q", url, http.StatusOK, resp.StatusCode, string(body),
		)
	}

	grandTotal, err := ParseGrandTotalResponse(body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse results %s", err)
	}
	return grandTotal, nil
}

func ParseGrandTotalResponse(data []byte) (*summary.GrandTotal, error) {
	var body summary.GrandTotal
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("Failed to parse json response: %s. body: %q", err, data)
	}
	return &body, nil
}

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
	TodaySummaryAPIRouter = "/api/summary/today"
)

func (c *Client) TodaySummary(ctx context.Context) (*summary.Summary, error) {
	url := c.baseURL + TodaySummaryAPIRouter
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Accept", "application/json")
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

	summary, err := ParseClientSummaryResponse(body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse today-summary results %s", err)
	}
	return summary, nil
}

func ParseClientSummaryResponse(data []byte) (*summary.Summary, error) {
	var body summary.Summary
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("Failed to parse json response: %s. body: %q", err, data)
	}
	return &body, nil
}

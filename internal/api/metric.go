package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/result17/codeBeatCli/internal/metric"
	"github.com/spf13/viper"
)

var metricKeyDataTypeMap = map[string]interface{}{
	"project": metric.MetricRatioData[string]{},
	"lineno":  metric.MetricRatioData[uint32]{},
}

var metricKeyParseFuncMap = map[string]func(data []byte) (interface{}, error){
	"project": func(data []byte) (interface{}, error) {
		return ParseStringMetricDurationResponse(data)
	},
	"lineno": func(data []byte) (interface{}, error) {
		return ParseIntMetricDurationResponse(data)
	},
}

func QueryTodayMetricDuration[T string | uint32](c *Client, ctx context.Context, v *viper.Viper) (*metric.MetricRatioData[T], error) {
	metricKey := v.GetString("today-metric-duration")
	url := fmt.Sprintf("%s/api/metric/duration/today/%s", c.baseURL, metricKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Fail to create request: %s", err)
	}

	resp, err := c.Do(ctx, req)

	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Fail to execute request: %s", err)
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

	if parseFunc, ok := metricKeyParseFuncMap[metricKey]; ok {
		data, err := parseFunc(body)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse today-metric-duration results %s", err)
		}
		if result, ok := data.(*metric.MetricRatioData[T]); ok {
			return result, nil
		}
		return nil, fmt.Errorf("Type mismatch for metric key %q", metricKey)
	}
	return nil, fmt.Errorf("Invalid metric key %q", metricKey)
}

func ParseStringMetricDurationResponse(data []byte) (*metric.MetricRatioData[string], error) {
	var body metric.MetricRatioData[string]
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("Failed to parse json response: %s. body: %q", err, data)
	}
	return &body, nil
}

func ParseIntMetricDurationResponse(data []byte) (*metric.MetricRatioData[uint32], error) {
	var body metric.MetricRatioData[uint32]
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("e: %s. body: %q", err, data)
	}
	return &body, nil
}

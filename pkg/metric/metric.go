package metric

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/result17/codeBeatCli/internal/api"
	"github.com/result17/codeBeatCli/internal/metric"
	apiCmd "github.com/result17/codeBeatCli/pkg/api"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

// Run executes the today-metric-duration command.
func Run[T string | uint32](ctx context.Context, v *viper.Viper) (int, error) {
	logger := log.Extract(ctx)
	data, err := TodayMetricDuration[T](ctx, v)
	if err != nil {
		logger.Errorf("Failed fetched today-summary, %s", err)
		return exitcode.ErrGeneric, fmt.Errorf(
			"Today fetch failed: %s",
			err,
		)
	}

	logger.Debugln("Successfully fetched today-summary")

	output, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Failed marshal summary, %s, %+v", err, data)
		return exitcode.ErrGeneric, fmt.Errorf(
			"Summary unmarshal failed: %s",
			err,
		)
	}
	fmt.Println(string(output))
	return exitcode.Success, nil
}

func TypeRun(ctx context.Context, v *viper.Viper, metricKey string) (int, error) {
	switch metricKey {
	case "project":
		return Run[string](ctx, v)
	case "lineno":
		return Run[uint32](ctx, v)
	default:
		return exitcode.ErrGeneric, fmt.Errorf("Invalid metric key: %s", metricKey)
	}
}

func TodayMetricDuration[T string | uint32](ctx context.Context, v *viper.Viper) (*metric.MetricRatioData[T], error) {
	apiParams, err := params.LoadApiParams(ctx, v)
	apiClient, err := apiCmd.NewClient(ctx, apiParams.BaseUrl)

	if err != nil {
		return nil, fmt.Errorf("Fail to create apiClient: %s", err)
	}
	data, err := api.QueryTodayMetricDuration[T](apiClient, ctx, v)
	if err != nil {
		return nil, fmt.Errorf("Fail to query today's summary: %s", err)
	}
	return data, nil
}

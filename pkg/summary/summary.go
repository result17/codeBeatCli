package summary

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/result17/codeBeatCli/internal/summary"
	apiCmd "github.com/result17/codeBeatCli/pkg/api"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

// Run executes the today-summary command.
func Run(ctx context.Context, v *viper.Viper) (int, error) {
	logger := log.Extract(ctx)
	summary, err := TodaySummary(ctx, v)
	if err != nil {
		logger.Errorf("Failed fetched today-summary, %s", err)
		return exitcode.ErrGeneric, fmt.Errorf(
			"Today fetch failed: %s",
			err,
		)
	}
	logger.Debugln("Successfully fetched today-summary")

	var output []byte

	err = json.Unmarshal(output, summary)
	if err != nil {
		logger.Errorf("Failed unmarshal summary, %s", err)
		return exitcode.ErrGeneric, fmt.Errorf(
			"Summary unmarshal failed: %s",
			err,
		)
	}
	fmt.Println(output)
	return exitcode.Success, nil
}

func TodaySummary(ctx context.Context, v *viper.Viper) (*summary.Summary, error) {
	apiParams, err := params.LoadApiParams(ctx, v)
	apiClient, err := apiCmd.NewClient(ctx, apiParams.BaseUrl)

	if err != nil {
		return nil, fmt.Errorf("Fail to create apiClient: %s", err)
	}
	summary, err := apiClient.TodaySummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("Fail to query today's summary: %s", err)
	}
	return summary, nil
}

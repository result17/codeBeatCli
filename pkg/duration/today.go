package duration

import (
	"context"
	"fmt"

	apiCmd "github.com/result17/codeBeatCli/pkg/api"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

// Run executes the today-duration command.
func Run(ctx context.Context, v *viper.Viper) (int, error) {
	logger := log.Extract(ctx)
	output, err := TodayDuration(ctx, v)
	if err != nil {
		logger.Errorf("Failed fetched today-duration for status bar, %s", err)
		return exitcode.ErrGeneric, fmt.Errorf(
			"Today fetch failed: %s",
			err,
		)
	}

	logger.Debugln("Successfully fetched today-duration for status bar")

	// stdout
	fmt.Println(output)

	return exitcode.Success, nil
}

func TodayDuration(ctx context.Context, v *viper.Viper) (string, error) {
	apiParams, err := params.LoadApiParams(ctx, v)
	apiClient, err := apiCmd.NewClient(ctx, apiParams.BaseUrl)

	if err != nil {
		return "", fmt.Errorf("Fail to create apiClient: %s", err)
	}

	grandTotal, err := apiClient.TodayDuration(ctx)
	if err != nil {
		return "", fmt.Errorf("Fail to query today's duration: %s", err)
	}
	// Returning text only so far
	return grandTotal.Text, nil
}

package today

import (
	"context"
	"fmt"

	apiCmd "github.com/result17/codeBeatCli/pkg/api"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

// Run executes the today command.
func Run(ctx context.Context, v *viper.Viper) (int, error) {
	output, err := Today(ctx, v)
	if err != nil {

		return exitcode.ErrGeneric, fmt.Errorf(
			"today fetch failed: %s",
			err,
		)
	}

	logger := log.Extract(ctx)

	logger.Debugln("successfully fetched today for status bar")

	// stdout
	fmt.Println(output)

	return exitcode.Success, nil
}

func Today(ctx context.Context, v *viper.Viper) (string, error) {
	apiParams, err := params.LoadApiParams(ctx, v)
	apiClient, err := apiCmd.NewClient(ctx, apiParams.BaseUrl)

	if err != nil {
		return "", fmt.Errorf("Fail to create apiClient: %s", err)
	}

	grandTotal, err := apiClient.Today(ctx)
	if err != nil {
		return "", fmt.Errorf("Fail to query today's duration: %s", err)
	}
	return grandTotal.Text, nil
}

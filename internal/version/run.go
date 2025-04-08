package version

import (
	"context"
	"fmt"

	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/spf13/viper"
)

func RunVersion(_ context.Context, v *viper.Viper) (int, error) {
	fmt.Printf(
		"CodeBeatCli\n Version: %s\n Commit: %s\n Built: %s\n OS/Arch: %s/%s\n",
		Version,
		Commit,
		BuildDate,
		OS,
		Arch,
	)
	return exitcode.Success, nil
}

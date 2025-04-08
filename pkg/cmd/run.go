package cmd

import (
	"bytes"
	"context"
	"io"
	stdlog "log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/result17/codeBeatCli/internal/version"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
)

// cmdFn represents a command function
type cmdFn func(ctx context.Context, v *viper.Viper) (int, error)

func RunE(cmd *cobra.Command, v *viper.Viper) error {
	ctx := context.Background()
	// add logger to context
	log.Extract(ctx)
	logger, err := SetupLogging(ctx)
	if err != nil {
		stdlog.Fatalf("failde to setup logging: %s", err)
	}
	ctx = log.ToContxt(ctx, logger)

	if v.GetBool("version") {
		logger.Debugln("command: version")
		return runCmd(ctx, v, version.RunVersion)
	}
	return exitcode.Err{Code: exitcode.ErrGeneric}
}

// TODO setup logger output file path
func SetupLogging(ctx context.Context) (*log.Logger, error) {
	var destOutput io.Writer = os.Stdout
	l := log.New(destOutput)
	return l, nil
}

func runCmd(ctx context.Context, v *viper.Viper, cmd cmdFn) (errorsponse error) {
	logs := bytes.NewBuffer(nil)
	resetLogs := captureLogs(ctx, logs)

	logger := log.Extract(ctx)

	var err error

	// run command
	exitCode, err := cmd(ctx, v)

	if err != nil {
		logger.Errorf("failed to run command: %s", err)

		resetLogs()
	}

	if exitCode != exitcode.Success {
		logger.Debugf("command failed with exit code %d", exitCode)
		errorsponse = exitcode.Err{Code: exitCode}
	}
	return errorsponse
}

func captureLogs(ctx context.Context, dest io.Writer) func() {
	logger := log.Extract(ctx)
	loggerOutput := logger.Output()

	mw := io.MultiWriter(loggerOutput, dest)

	logger.SetOutput(mw)

	return func() {
		logger.SetOutput(loggerOutput)
	}
}

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	"github.com/result17/codeBeatCli/internal/version"
	"github.com/result17/codeBeatCli/pkg/duration"
	heartbeat "github.com/result17/codeBeatCli/pkg/entity"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/summary"
	"gopkg.in/natefinch/lumberjack.v2"
)

// cmdFn represents a command function
type cmdFn func(ctx context.Context, v *viper.Viper) (int, error)

func RunE(cmd *cobra.Command, v *viper.Viper) error {
	ctx := context.Background()
	// add logger to context
	log.Extract(ctx)
	logger, err := SetupLogging(ctx, v)
	if err != nil {
		stdlog.Fatalf("Failde to setup logging: %s", err)
	}
	ctx = log.ToContxt(ctx, logger)

	if v.GetBool("version") {
		logger.Debugln("command: version")
		return runCmd(ctx, v, version.RunVersion)
	}

	if entity := v.GetString("entity"); entity != "" {
		logger.Debugln("Command: heartbeat")
		_, err := heartbeat.Run(ctx, v)
		return err
	}

	if v.GetBool("today-duration") {
		logger.Debugln("Command: today-duration")
		_, err := duration.Run(ctx, v)
		return err
	}

	if v.GetBool("today-summary") {
		logger.Debugln("command: today-summary")
		_, err := summary.Run(ctx, v)
		return err
	}

	_ = cmd.Help()
	return exitcode.Err{Code: exitcode.ErrGeneric}
}

func SetupLogging(ctx context.Context, v *viper.Viper) (*log.Logger, error) {
	var destOutput io.Writer = os.Stdout

	logPath, err := log.LogFilepath()

	if err != nil {
		return nil, fmt.Errorf("Fail to output log file %s", logPath)
	}

	destOutput = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    log.MaxLogFileSize,
		MaxBackups: log.MaxNumberOfBackups,
	}
	l := log.New(destOutput)
	if v.GetBool("dlog") {
		l.SetAtomicLevel(zapcore.DebugLevel)
	}
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
		logger.Errorf("Failed to run command: %s", err)

		resetLogs()
	}

	if exitCode != exitcode.Success {
		logger.Debugf("Command failed with exit code %d", exitCode)
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

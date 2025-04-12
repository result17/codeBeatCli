package heartbeat

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/matishsiao/goInfo"
	"github.com/result17/codeBeatCli/internal/heartbeat"
	"github.com/result17/codeBeatCli/internal/offline"
	"github.com/result17/codeBeatCli/internal/version"
	apiCmd "github.com/result17/codeBeatCli/pkg/api"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

func SendHeartbeats(ctx context.Context, v *viper.Viper, path string) error {
	h, err := params.LoadHeartbeatParams(ctx, v)
	if err != nil {
		return fmt.Errorf("failed to load heartbeat parameters: %w", err)
	}

	logger := log.Extract(ctx)
	setLogFields(logger, h)
	logger.Debugf("params: %s", h)

	opts := initHandleOptions(h)
	if isSave := v.GetBool("local-save"); isSave {
		opts = append(opts, offline.SaveHeartbeat(path))
	}
	heartbeats := buildHeartbeats(ctx, h)
	// TODO RateLimit
	// TODO backoff handler

	apiClient, err := apiCmd.NewClient(ctx)

	handle := heartbeat.NewHandle(apiClient, opts...)
	results, err := handle(ctx, heartbeats)

	for _, result := range results {
		if len(result.Errors) > 0 {
			logger.Warnln(strings.Join(result.Errors, " "))
		}
	}

	return nil
}

func setLogFields(logger *log.Logger, params params.Heartbeat) {
	logger.AddField("entity", params.Entity)
	logger.AddField("time", params.Time)

	if params.LineNumber != nil {
		logger.AddField("lineno", params.LineNumber)
	}
}

func UserAgent(ctx context.Context, plugin string) string {
	logger := log.Extract(ctx)

	template := "codeBeat/%s (%s-%s-%s) %s %s"

	if plugin == "" {
		plugin = "codeBeat-v0/"
	}

	info, err := goInfo.GetInfo()
	if err != nil {
		logger.Debugf("goInfo.GetInfo error: %s", err)
	}

	userAgent := fmt.Sprintf(
		template,
		version.Version,
		// TODO system pkg
		"windows",
		strings.TrimSpace(info.Core),
		strings.TrimSpace(info.Platform),
		strings.TrimSpace(runtime.Version()),
		strings.TrimSpace(plugin),
	)

	defer func() {
		if r := recover(); r != nil {
			userAgent = fmt.Sprintf(
				template,
				version.Version,
				// TODO system pkg
				"windows",
				"unknown",
				"unknown",
				strings.TrimSpace(runtime.Version()),
				strings.TrimSpace(plugin),
			)
		}
	}()
	return userAgent
}

func buildHeartbeats(ctx context.Context, params params.Heartbeat) []heartbeat.Heartbeat {
	heartbeats := []heartbeat.Heartbeat{}
	userAgent := UserAgent(ctx, params.Plugin)

	heartbeats = append(heartbeats, *heartbeat.New(
		params.Entity,
		userAgent,
		params.Time,
		params.CursorPos,
		params.Language,
		params.LineNumber,
		params.LineInFile,
		&params.AlternateProject,
		&params.ProjectFolder,
	))
	return heartbeats
}

func initHandleOptions(params params.Heartbeat) []heartbeat.HandleOption {
	opts := []heartbeat.HandleOption{
		heartbeat.WithFormatting(),
	}
	return opts
}

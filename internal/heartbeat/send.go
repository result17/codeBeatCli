package heartbeat

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/matishsiao/goInfo"
	"github.com/result17/codeBeatCli/internal/version"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/result17/codeBeatCli/pkg/params"
	"github.com/spf13/viper"
)

func SendHeartbeats(ctx context.Context, v *viper.Viper, path string) error {
	heartbeat, err := params.LoadHeartbeatParams(ctx, v)
	if err != nil {
		return fmt.Errorf("failed to load heartbeat parameters: %w", err)
	}

	logger := log.Extract(ctx)
	setLogFields(logger, heartbeat)
	logger.Debugf("params: %s", heartbeat)
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

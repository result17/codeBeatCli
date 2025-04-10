package heartbeat

import (
	"context"
	"path/filepath"

	"github.com/result17/codeBeatCli/internal/windows"
	"github.com/result17/codeBeatCli/pkg/log"
)

func WithFormatting() HandleOption {
	return func(next Handle) Handle {
		return func(ctx context.Context, hs []Heartbeat) ([]Result, error) {
			logger := log.Extract(ctx)
			logger.Debugln("execute heartbeat filepath formatting")

			for n, h := range hs {
				hs[n] = Format(ctx, h)
			}
			return next(ctx, hs)
		}
	}
}

func Format(ctx context.Context, h Heartbeat) Heartbeat {
	return h
}

func formatWindowsFilePath(ctx context.Context, h *Heartbeat) {
	logger := log.Extract(ctx)

	formatted, err := filepath.Abs(h.Entity)
	if err != nil {
		logger.Debugf("failed to resolve absolute path for %q: %s", h.Entity, err)
		return
	}

	h.Entity = windows.FormatFilePath(formatted)
}

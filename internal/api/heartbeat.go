package api

import (
	"context"
	"errors"

	"github.com/result17/codeBeatCli/internal/heartbeat"
	"github.com/result17/codeBeatCli/pkg/log"
)

func (c Client) SendHeartbeats(ctx context.Context, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
	logger := log.Extract(ctx)
	logger.Debugf("Sending %d heartbeats(s) to api at %s", len(hs), c.baseURL)

	return nil, errors.New("TODO: sending heartbeat failed")
}

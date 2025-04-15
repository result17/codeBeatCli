package heartbeat

import (
	"context"

	"github.com/result17/codeBeatCli/internal/offline"
	"github.com/result17/codeBeatCli/pkg/exitcode"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/spf13/viper"
)

func Run(ctx context.Context, v *viper.Viper) (int, error) {
	logger := log.Extract(ctx)
	queueFilepath, err := offline.QueueFilepath(ctx, v)
	if err != nil {
		logger.Warnf("Fail to load offline queue filepath: %s", err)
	}

	err = SendHeartbeats(ctx, v, queueFilepath)
	if err != nil {
		logger.Debugln("Fail to sent heartbeat(s)")
		return exitcode.ErrAPI, nil
	}

	logger.Debugln("successfully sent heartbeat(s)")

	return exitcode.Success, nil
}

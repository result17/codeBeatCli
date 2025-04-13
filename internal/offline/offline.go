package offline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/result17/codeBeatCli/internal/heartbeat"
	"github.com/result17/codeBeatCli/internal/workspace"
	"github.com/result17/codeBeatCli/pkg/log"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

const (
	// dbFilename is the default bolt db filename.
	dbFilename = "offline_heartbeats_codebeat.bdb"
	// maxRequeueAttempts defines the maximum number of attempts to requeue heartbeats,
	// which could not successfully be sent to the WakaTime API.
	maxRequeueAttempts = 3
)

func openDB(ctx context.Context, fp string) (db *bolt.DB, _ func(), err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("OpenDB panicked")
		}
	}()
	logger := log.Extract(ctx)
	logger.Debugf("Open db file: %s", fp)
	db, err = bolt.Open(fp, 0644, &bolt.Options{Timeout: 30 * time.Second})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open db file: %s", err)
	}

	return db, func() {
		logger := log.Extract(ctx)

		// recover from panic when closing db
		defer func() {
			if r := recover(); r != nil {
				logger.Warnf("panicked: failed to close db file: %v", r)
			}
		}()

		if err := db.Close(); err != nil {
			logger.Debugf("failed to close db file: %s", err)
		}
	}, err
}

func QueueFilepath(ctx context.Context, v *viper.Viper) (string, error) {
	homedir, err := workspace.CodeBeatHomeDir()

	if err != nil {
		return dbFilename, fmt.Errorf("failed getting resource directory, defaulting to current directory: %s", err)
	}
	return filepath.Join(homedir, ".codebeat", dbFilename), nil
}

func WithQueue(fp string) heartbeat.HandleOption {
	return func(next heartbeat.Handle) heartbeat.Handle {
		return func(ctx context.Context, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
			logger := log.Extract(ctx)
			logger.Debugf("execute offline queue with file %s", fp)

			if len(hs) == 0 {
				logger.Debugln("abort execution, as there are no heartbeats ready for sending")

				return nil, nil
			}
			results, err := next(ctx, hs)
			if err != nil {
				logger.Debugf("pushing %d heartbeat(s) to queue after error: %s", len(hs), err)

				requeueErr := pushHeartbeatsWithRetry(ctx, fp, hs)
				if requeueErr != nil {
					return nil, fmt.Errorf(
						"failed to push heartbeats to queue: %s",
						requeueErr,
					)
				}
				return nil, err
			}
			err = handleResults(ctx, fp, results, hs)
			if err != nil {
				return nil, fmt.Errorf("failed to handle results: %s", err)
			}
			return results, nil
		}
	}
}

func SaveHeartbeat(fp string) heartbeat.HandleOption {
	return func(next heartbeat.Handle) heartbeat.Handle {
		return func(ctx context.Context, hs []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
			requeueErr := pushHeartbeatsWithRetry(ctx, fp, hs)
			if requeueErr != nil {
				return nil, fmt.Errorf(
					"saving heartbeat locally failed to push heartbeats to queue: %s",
					requeueErr,
				)
			}
			results, err := next(ctx, hs)
			return results, err
		}
	}
}

func pushHeartbeatsWithRetry(ctx context.Context, fp string, hs []heartbeat.Heartbeat) error {
	var (
		count int
		err   error
	)

	logger := log.Extract(ctx)

	for {
		if count >= maxRequeueAttempts {
			serialized, jsonErr := json.Marshal(hs)
			if jsonErr != nil {
				logger.Warnf("failed to json marshal heartbeats: %s. heartbeats: %#v", jsonErr, hs)
			}

			return fmt.Errorf(
				"abort requeuing after %d unsuccessful attempts: %s. heartbeats: %s",
				count,
				err,
				string(serialized),
			)
		}
		err = pushHeartbeats(ctx, fp, hs)
		if err != nil {
			count++
			sleepSeconds := math.Pow(2, float64(count))

			time.Sleep(time.Duration(sleepSeconds) * time.Second)

			continue
		}
		break
	}
	return nil
}

func pushHeartbeats(ctx context.Context, fp string, hs []heartbeat.Heartbeat) error {
	db, close, err := openDB(ctx, fp)
	if err != nil {
		return err
	}
	defer close()

	tx, err := db.Begin(true)
	if err != nil {
		return fmt.Errorf("failed to start db transaction: %s", err)
	}

	queue := NewQueue(tx)
	err = queue.PushMany(hs)
	if err != nil {
		return fmt.Errorf("failed to push heartbeat(s) to queue: %s", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit db transaction: %s", err)
	}
	return nil
}

// TODO handle API response
func handleResults(ctx context.Context, fp string, results []heartbeat.Result, hs []heartbeat.Heartbeat) error {
	return nil
}

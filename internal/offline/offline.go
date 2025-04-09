package offline

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/result17/codeBeatCli/internal/workspace"
	"github.com/spf13/viper"
)

const (
	// dbFilename is the default bolt db filename.
	dbFilename = "offline_heartbeats.bdb"
)

func QueueFilepath(ctx context.Context, v *viper.Viper) (string, error) {
	homedir, err := workspace.CodeBeatHomeDir()
	if err != nil {
		return dbFilename, fmt.Errorf("failed getting resource directory, defaulting to current directory: %s", err)
	}
	return filepath.Join(homedir, dbFilename), nil
}

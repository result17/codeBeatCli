package workspace

import (
	"os"
)

func CodeBeatHomeDir() (string, error) {
	return os.UserHomeDir()
}

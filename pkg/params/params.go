package params

import (
	"context"
	"errors"
	"time"

	"github.com/result17/codeBeatCli/internal/vipertools"
	"github.com/spf13/viper"
)

type (
	Heartbeat struct {
		Entity           string
		Plugin           string
		Language         *string
		LineNumber       *int
		CursorPos        *int
		LineInFile       *int
		AlternateProject string
		ProjectFolder    string
		Config           *string
		LogFile          *string
		Time             int64
	}
)

// PointerTo returns a pointer to the value passed in.
func PointerTo[t bool | int | string](v t) *t {
	return &v
}

func LoadHeartbeatParams(ctx context.Context, v *viper.Viper) (Heartbeat, error) {
	var cursorPos *int
	if v.IsSet("cursorpos") {
		pos := v.GetInt("cursorpos")
		cursorPos = PointerTo(pos)
	}
	entity := vipertools.GetString(v, "entity")
	if entity == "" {
		return Heartbeat{}, errors.New("fail to receive entity")
	}

	var lineNumber *int
	if v.IsSet("lineno") {
		lineNumber = PointerTo(v.GetInt("lineno"))
	}

	var lineInFile *int
	if v.IsSet("line-in-file") {
		lineInFile = PointerTo(v.GetInt("line-in-file"))
	}

	plugin := vipertools.GetString(v, "plugin")
	alternateProject := vipertools.GetString(v, "alternate-project")
	projectFloader := vipertools.GetString(v, "project-floader")

	var language *string
	if l := vipertools.GetString(v, "language"); l != "" {
		language = &l
	}

	// default now
	timeVal := time.Now().Unix()
	if v.IsSet("time") {
		if secs := v.GetInt64("time"); secs > 0 {
			timeVal = secs
		}
	}

	return Heartbeat{
		Entity:           entity,
		Plugin:           plugin,
		LineNumber:       lineNumber,
		CursorPos:        cursorPos,
		LineInFile:       lineInFile,
		AlternateProject: alternateProject,
		ProjectFolder:    projectFloader,
		Time:             timeVal,
		Language:         language,
	}, nil
}

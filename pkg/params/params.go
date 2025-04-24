package params

import (
	"context"
	"errors"
	"time"

	"github.com/result17/codeBeatCli/internal/api"
	"github.com/result17/codeBeatCli/internal/vipertools"
	"github.com/spf13/viper"
)

type (
	Params struct {
		API       API
		Heartbeat Heartbeat
	}

	API struct {
		BaseUrl string
	}

	Heartbeat struct {
		Entity           string
		Plugin           string
		Language         *string
		LinesNumber      *int
		CursorPos        *int
		LineInFile       *int
		AlternateProject *string
		ProjectFolder    *string
		Config           *string
		LogFile          *string
		Time             int64
	}
)

// PointerTo returns a pointer to the value passed in.
func PointerTo[t bool | int | string](v t) *t {
	return &v
}

func LoadParams(ctx context.Context, v *viper.Viper) (Params, error) {
	apiParams, err := loadApiParams(ctx, v)
	if err != nil {
		return Params{}, err
	}
	heartbeatParams, err := loadHeartbeatParams(ctx, v)
	if err != nil {
		return Params{}, err
	}
	return Params{
		API:       apiParams,
		Heartbeat: heartbeatParams,
	}, nil
}

func loadApiParams(ctx context.Context, v *viper.Viper) (API, error) {
	var baseUrl string
	if baseUrl = vipertools.GetString(v, "api-url"); baseUrl == "" {
		baseUrl = api.BaseURL
	}
	return API{
		BaseUrl: baseUrl,
	}, nil
}

func loadHeartbeatParams(ctx context.Context, v *viper.Viper) (Heartbeat, error) {
	var cursorPos *int
	if v.IsSet("cursorpos") {
		pos := v.GetInt("cursorpos")
		cursorPos = PointerTo(pos)
	}
	entity := vipertools.GetString(v, "entity")
	if entity == "" {
		return Heartbeat{}, errors.New("fail to receive entity")
	}

	var linesNumber *int
	if v.IsSet("lineno") {
		linesNumber = PointerTo(v.GetInt("lineno"))
	}

	var lineInFile *int
	if v.IsSet("lines-in-file") {
		lineInFile = PointerTo(v.GetInt("lines-in-file"))
	}

	plugin := vipertools.GetString(v, "plugin")

	var alternateProject *string
	if v.IsSet("alternate-project") {
		alternateProject = PointerTo(vipertools.GetString(v, "alternate-project"))
	}

	var projectFloder *string
	if v.IsSet("project-path") {
		projectFloder = PointerTo(vipertools.GetString(v, "project-path"))
	}

	var language *string
	if l := vipertools.GetString(v, "language"); l != "" {
		language = &l
	}

	// default now
	timeVal := time.Now().UnixMilli()
	if v.IsSet("time") {
		if secs := v.GetInt64("time"); secs > 0 {
			timeVal = secs
		}
	}

	return Heartbeat{
		Entity:           entity,
		Plugin:           plugin,
		LinesNumber:      linesNumber,
		CursorPos:        cursorPos,
		LineInFile:       lineInFile,
		AlternateProject: alternateProject,
		ProjectFolder:    projectFloder,
		Time:             timeVal,
		Language:         language,
	}, nil
}

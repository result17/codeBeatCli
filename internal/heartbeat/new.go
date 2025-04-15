package heartbeat

import (
	"context"
	"fmt"
)

type Result struct {
	Errors    []string
	Status    int
	Heartbeat Heartbeat
}

type Sender interface {
	SendHeartbeats(context.Context, []Heartbeat) ([]Result, error)
}

// Handle does processing of heartbeats.
type Handle func(context.Context, []Heartbeat) ([]Result, error)

// HandleOption is a function, which allows chaining multiple Handles.
type HandleOption func(next Handle) Handle

func NewHandle(
	sender Sender, opts ...HandleOption) Handle {
	return func(ctx context.Context, hs []Heartbeat) ([]Result, error) {
		var handle Handle = sender.SendHeartbeats
		for i := len(opts) - 1; i >= 0; i-- {
			handle = opts[i](handle)
		}
		return handle(ctx, hs)
	}
}

// TODO heartbeat ID
type Heartbeat struct {
	CursorPosition *int    `json:"cursorpos,omitempty"`
	Entity         string  `json:"entity"`
	Language       *string `json:"language,omitempty"`
	LineNumber     *int    `json:"lineno,omitempty"`
	LinesInFile    *int    `json:"lines,omitempty"`
	Project        *string `json:"project,omitempty"`
	ProjectPath    *string `json:"projectPath,omitempty"`
	Time           int64   `json:"time"`
	UserAgent      string  `json:"user_agent"`
}

func New(entity, userAgent string, time int64, cursorPos *int, lang *string, lineNum *int, linesInFile *int, project *string, projectPath *string) *Heartbeat {
	hb := &Heartbeat{
		Entity:         entity,
		Time:           time,
		UserAgent:      userAgent,
		CursorPosition: cursorPos,
		Language:       lang,
		LineNumber:     lineNum,
		LinesInFile:    linesInFile,
		Project:        project,
		ProjectPath:    projectPath,
	}

	return hb
}

func (h Heartbeat) ID() string {
	project := "unset"
	if h.Project != nil {
		project = *h.Project
	}

	cursorPos := "nil"
	if h.CursorPosition != nil {
		cursorPos = fmt.Sprint(*h.CursorPosition)
	}

	return fmt.Sprintf("%d-%s-%s", h.Time, cursorPos, project)
}

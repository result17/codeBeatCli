package heartbeat

type Heartbeat struct {
	CursorPosition *int    `json:"cursorpos,omitempty"`
	Entity         string  `json:"entity"`
	Language       *string `json:"language,omitempty"`
	LineNumber     *int    `json:"lineno,omitempty"`
	LinesInFile    *int    `json:"lines,omitempty"`
	Project        *string `json:"project,omitempty"`
	ProjectPath    *string `json:"project,omitempty"`
	Time           int64   `json:"time"`
	UserAgent      string  `json:"user_agent"`
}

func New(entity, projectPath, userAgent string, time int64, cursorPos *int, lang *string, lineNum *int, linesInFile *int, project *string, projectPath *string) *Heartbeat {
	hb := &Heartbeat{
		Entity:      entity,
		ProjectPath: projectPath,
		Time:        time,
		UserAgent:   userAgent,
		CursorPosition: cursorPos,
		Language: lang,
		LineNumber: lineNum,
		LinesInFile: linesInFile,
		Project: project,
		ProjectPath: projectPath
	}

	return hb
}

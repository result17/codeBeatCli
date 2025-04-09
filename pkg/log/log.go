package log

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/result17/codeBeatCli/internal/version"
)

type Logger struct {
	entry              *zap.Logger
	atomicLever        zap.AtomicLevel
	output             io.Writer
	dynamicWriteSyncer *DynamicWriteSyncer
}

func New(dest io.Writer, opts ...Option) *Logger {
	// default is info level
	level := zap.NewAtomicLevel()
	writer := NewDynamicWriteSyncer(zapcore.AddSync(dest))

	encordCfg := zap.NewProductionEncoderConfig()
	encordCfg.FunctionKey = "func"
	encordCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encordCfg),
		writer,
		level,
	), zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))

	logger = logger.With(
		zap.String("version", version.Version),
		zap.String("os/arch", fmt.Sprintf("%s/%s", version.OS, version.Arch)),
	)

	l := &Logger{
		entry:              logger,
		atomicLever:        level,
		output:             dest,
		dynamicWriteSyncer: writer,
	}

	for _, option := range opts {
		option(l)
	}

	return l
}

func (l *Logger) Infof(format string, args ...any) {
	l.entry.Log(zap.InfoLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.entry.Log(zap.DebugLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) Debugln(msg string) {
	l.entry.Log(zap.DebugLevel, msg)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.entry.Log(zapcore.ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.entry.Log(zapcore.WarnLevel, fmt.Sprintf(format, args...))
}

func (l *Logger) Output() io.Writer {
	return l.output
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
	l.dynamicWriteSyncer.SetWriter(zapcore.AddSync(w))
}

// SetAtomicLevel temporarily sets the logger's level and returns a restore function
// that will revert to the previous level when called.
// This is useful for temporarily changing log levels in specific code sections.
func (l *Logger) SetAtomicLevel(level zapcore.Level) (restore func()) {
	prevLevel := l.atomicLever.Level()
	l.atomicLever.SetLevel(level)

	// Named return value makes it clear what the function returns
	restore = func() {
		l.atomicLever.SetLevel(prevLevel)
	}
	return
}

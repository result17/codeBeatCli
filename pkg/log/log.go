package log

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	entry              *zap.Logger
	output             io.Writer
	dynamicWriteSyncer *DynamicWriteSyncer
}

func New(dest io.Writer) *Logger {
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

	l := &Logger{
		entry:              logger,
		output:             dest,
		dynamicWriteSyncer: writer,
	}

	return l
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

func (l *Logger) Output() io.Writer {
	return l.output
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
	l.dynamicWriteSyncer.SetWriter(zapcore.AddSync(w))
}

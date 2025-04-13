package log

import (
	"context"
	"os"
)

type (
	ctxSignal struct{}
)

var ctxSignalKey = &ctxSignal{}

func Extract(ctx context.Context) *Logger {
	l, ok := ctx.Value(ctxSignalKey).(*Logger)
	if !ok || l == nil {
		return New(os.Stdout)
	}
	return l
}

func ToContxt(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, ctxSignalKey, l)
}

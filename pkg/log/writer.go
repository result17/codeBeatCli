package log

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

type DynamicWriteSyncer struct {
	mu     sync.RWMutex
	writer zapcore.WriteSyncer
}

func (d *DynamicWriteSyncer) Write(p []byte) (n int, err error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.writer.Write(p)
}

func (d *DynamicWriteSyncer) Sync() error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.writer.Sync()
}

func (d *DynamicWriteSyncer) SetWriter(writer zapcore.WriteSyncer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.writer = writer
}

func NewDynamicWriteSyncer(inital zapcore.WriteSyncer) *DynamicWriteSyncer {
	return &DynamicWriteSyncer{
		writer: inital,
	}
}

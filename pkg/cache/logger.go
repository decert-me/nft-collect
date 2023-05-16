package cache

import (
	"fmt"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Zap *zap.Logger
}

func NewZapLogger(zap *zap.Logger) *ZapLogger {
	return &ZapLogger{zap}
}

func (c *ZapLogger) Printf(format string, v ...interface{}) {
	c.Zap.Log(2, fmt.Sprintf(format, v...))
}

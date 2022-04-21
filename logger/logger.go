package logger

import (
	"context"
	"lib/config"
	"sync"

	"go.uber.org/zap"
)

var (
	sugarLogger *zap.SugaredLogger
	syncMutex   sync.Mutex
)

// InitLogger 初始化日志
func InitLogger(c config.ZapLoggerConfig) error {
	if sugarLogger == nil {
		syncMutex.Lock()
		defer syncMutex.Unlock()
		logger, err := newZapLogger(c)
		if err != nil {
			return err
		}
		sugarLogger = logger
	}
	return nil
}

// WithLogger 根据上下文拼接日志
func WithLogger(ctx context.Context, withs ...interface{}) *zap.SugaredLogger {
	return FromLoggerContext(ctx).With(withs...)
}

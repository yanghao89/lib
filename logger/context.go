package logger

import (
	"context"
	"go.uber.org/zap"
)

type (
	loggerCtxKey struct{}
)

// NewLoggerContext 初始化上下文
func NewLoggerContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

// FromLoggerContext 从上下文中获取日志结构体
func FromLoggerContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerCtxKey{}).(*zap.SugaredLogger); ok {
		return logger
	}
	return Logger()
}

// Logger 获取日志
func Logger() *zap.SugaredLogger {
	if sugarLogger == nil {
		panic("sugaredLogger is nil")
	}
	return sugarLogger
}

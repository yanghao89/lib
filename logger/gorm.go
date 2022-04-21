package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	glog "gorm.io/gorm/logger"
)

type (
	gorm glog.Interface
)

// InitGormLogger 初始化 gorm 日志
func InitGormLogger(level int) glog.Interface {
	return newGormLogger().LogMode(glog.LogLevel(level))
}

var (
	_ gorm = (*GormLogger)(nil)
)

type GormLogger struct {
	LogLevel                  glog.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func newGormLogger() *GormLogger {
	return &GormLogger{
		LogLevel:                  glog.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          true,
		IgnoreRecordNotFoundError: false,
	}
}

func (g *GormLogger) LogMode(l glog.LogLevel) glog.Interface {
	g.LogLevel = l
	return g
}

func (g *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if g.LogLevel < glog.Info {
		return
	}
	//输出日志
	FromLoggerContext(ctx).Debugf(str, args...)
}

func (g *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if g.LogLevel < glog.Warn {
		return
	}
	//输出日志
	FromLoggerContext(ctx).Warnf(str, args...)
}

func (g *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if g.LogLevel < glog.Error {
		return
	}
	//输出日志
	FromLoggerContext(ctx).Warnf(str, args...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fun func() (string, int64), err error) {
	if g.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fun()
	var (
		rowInterface interface{} = rows
	)
	if rows == -1 {
		rowInterface = "-"
	}
	//从中间件获取日志对象
	logger := FromLoggerContext(ctx)
	ms := float64(elapsed.Nanoseconds()) / 1e6
	switch {
	case err != nil && g.LogLevel >= glog.Error && (!g.IgnoreRecordNotFoundError || !errors.Is(err, glog.ErrRecordNotFound)):
		logger.Errorf("\n[%.3fms] [rows:%v] %s",
			ms,
			rowInterface,
			sql,
		)
	case g.LogLevel == glog.Info:
		logger.Infof("\n[%.3fms] [rows:%v] %s",
			ms,
			rowInterface,
			sql,
		)
	case g.SlowThreshold != 0 && elapsed > g.SlowThreshold && g.LogLevel >= glog.Warn:
		logger.Warnf("%v\n[%.3fms] [rows:%v] %s",
			zap.Duration("SLOW SQL >= %v", g.SlowThreshold),
			ms,
			rowInterface,
			sql,
		)
	}
}

package logger

import (
	"io"
	"lib/config"
	"os"
	"path/filepath"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func writerFile(c config.ZapLoggerConfig) (io.Writer, error) {
	filePath := filepath.Join(c.FilePath, time.Now().Format("200601"))
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filePath, 0755)
		} else if os.IsPermission(err) {
			return nil, err
		}
	}
	return rotate.New(
		filePath+"/"+c.FileName+".%Y%m%d.log",
		rotate.WithLinkName(c.FilePath+c.FileName+".log"),
		rotate.WithMaxAge(time.Hour*24*time.Duration(c.MaxAge)),
		rotate.WithRotationTime(time.Hour*24),
	)
}

func newZapLogger(c config.ZapLoggerConfig) (*zap.SugaredLogger, error) {
	writer, err := writerFile(c)
	if err != nil {
		return nil, err
	}
	var (
		sync zapcore.WriteSyncer
	)
	if c.ShowConsole {
		sync = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(writer),
		)
	} else {
		sync = zapcore.AddSync(writer)
	}
	//自定义时间出书格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "line",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      customLevelEncoder, // 小写编码器
		EncodeTime:       customTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder, //
		EncodeCaller:     customCallerEncoder,            // 全路径编码器
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " | ",
	}

	var (
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		level   = getLoggerLevel(c.Level)
	)
	core := zapcore.NewCore(
		encoder,
		sync,
		level,
	)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar(), nil
}

func getLoggerLevel(l string) zapcore.Level {
	level := zap.InfoLevel
	switch l {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	return level
}

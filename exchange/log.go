package exchange

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
)

var LogSys *zap.Logger
var LogErr *zap.Logger

func InitLog() {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/home/logs/integrator/system/log-system.log",
		MaxSize:    500, // megabytes
		MaxBackups: 20,
		MaxAge:     14, // days
	})

	we := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/home/logs/integrator/error/log-error.log",
		MaxSize:    500, // megabytes
		MaxBackups: 20,
		MaxAge:     14, // days
	})

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	coreSys := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		ws,
		zap.InfoLevel,
	)

	coreErr := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		we,
		zap.InfoLevel,
	)
	LogSys = zap.New(coreSys)
	LogErr = zap.New(coreErr)
}

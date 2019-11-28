package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	Logger = initLogger("./order.log", zapcore.InfoLevel).Sugar()
}

func initLogger(path string, level zapcore.Level) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	atom := zap.NewAtomicLevelAt(level)

	config := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		InitialFields:    map[string]interface{}{"serviceName": "sorter_ledger"},
		OutputPaths:      []string{"stdout", path},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build log.
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}
	return logger
}

package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLog(t *testing.T) {
	err := InitLogger("../order.log", zapcore.InfoLevel)
	assert.NoError(t, err)

	Logger.Info("test logger")
	Logger.Infof("test %s", "logger")
	Logger.Infow("test logger", "k", "v")
}

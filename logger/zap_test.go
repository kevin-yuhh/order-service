package logger

import (
	"testing"
)

func TestLog(t *testing.T) {
	Logger.Info("test logger")
	Logger.Infof("test %s", "logger")
	Logger.Infow("test logger", "k", "v")
}

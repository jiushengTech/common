package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	for range 10000 {
		time.Sleep(1 * time.Second)
		Slog.Debug("test debug")
		Slog.Debugf("test debugf %s", "debugf")
		Slog.Debugw("test debugw", "debugw")
		Slog.Info("test info")
		Slog.Infof("test infof %s", "infof")
		Slog.Infow("test infow", "infow")
		Slog.Warn("test warn")
		Slog.Warnf("test warnf %s", "warnf")
		Slog.Warnw("test warnw", "warnw")
		Slog.Error("test error")
		Slog.Errorf("test errorf %s", "errorf")
		Slog.Errorw("test errorw", "errorw")

	}
}

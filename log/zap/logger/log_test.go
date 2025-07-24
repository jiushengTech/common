package logger

import (
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	r := func() {
		for range 100 {
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
	group := sync.WaitGroup{}
	group.Add(10)
	for _ = range 10 {
		go func() {
			defer group.Done()
			r()
		}()
	}
	group.Wait()

}

package logger

import (
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	r := func() {
		for range 100 {
			Log.Debug("test debug")
			Slog.Debugf("test debugf %s", "debugf")
			Log.Info("test info")
			Slog.Infof("test infof %s", "infof")
			Log.Warn("test warn")
			Slog.Warnf("test warnf %s", "warnf")
			Log.Error("test error")
			Slog.Errorf("test errorf %s", "errorf")
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

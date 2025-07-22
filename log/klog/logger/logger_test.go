package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Init()
	for range 10000 {
		time.Sleep(1 * time.Second)
		Log.Debug("test debug")
		Log.Debugf("test debugf %s", "debugf")
		Log.Debugw("test debugw", "debugw")
		Log.Info("test info")
		Log.Infof("test infof %s", "infof")
		Log.Infow("test infow", "infow")
		Log.Warn("test warn")
		Log.Warnf("test warnf %s", "warnf")
		Log.Warnw("test warnw", "warnw")
		Log.Error("test error")
		Log.Errorf("test errorf %s", "errorf")
		Log.Errorw("test errorw", "errorw")

	}
}

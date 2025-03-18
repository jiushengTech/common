package zap

import "testing"

func TestLog(t *testing.T) {
	for range 100000 {
		log.Debug("test debug")
		log.Debugf("test debugf %s", "debugf")
		log.Debugw("test debugw", "debugw", "debugw")
		log.Info("test info")
		log.Infof("test infof %s", "infof")
		log.Infow("test infow", "infow", "infow")
		log.Warn("test warn")
		log.Warnf("test warnf %s", "warnf")
		log.Warnw("test warnw", "warnw", "warnw")
		log.Error("test error")
		log.Errorf("test errorf %s", "errorf")
		log.Errorw("test errorw", "errorw", "errorw")
	}

}

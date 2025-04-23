package logger

import (
	"testing"
)

func TestLog(t *testing.T) {

	for range 10000000 {
		//time.Sleep(1 * time.Second)
		Slog.Debug("test debug")
	}
}

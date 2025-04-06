package logger

import (
	z "github.com/jiushengTech/common/log/zap"
	"go.uber.org/zap"
)

var Log *zap.Logger
var Slog *zap.SugaredLogger

func init() {
	Log = z.DefaultZapLogger()
	Slog = Log.Sugar()
}

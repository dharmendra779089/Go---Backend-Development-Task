package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	Log = logger
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

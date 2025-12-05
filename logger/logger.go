package logger

import (
	"cngamesdk.com/cron-task/global"
	"go.uber.org/zap"
)

type CronLog struct {
}

func (receiver CronLog) Info(msg string, keysAndValues ...interface{}) {
	global.Logger.Info("cron-info", zap.String("msg", msg), zap.Any("keysAndValues", keysAndValues))
}

func (receiver CronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	global.Logger.Error("cron-error", zap.Error(err), zap.String("msg", msg), zap.Any("keysAndValues", keysAndValues))
}

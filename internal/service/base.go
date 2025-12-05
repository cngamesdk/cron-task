package service

import (
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"context"
)

type BaseService struct {
	IsRunning bool
	Config    *cron_task.DimCronTaskConfigModel
	TaskLog   *cron_task.OdsCronTaskLogModel
}

func (receiver *BaseService) CompleteEvent(ctx context.Context) (err error) {
	return receiver.TaskLog.Create(ctx)
}

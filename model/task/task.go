package task

import (
	"cngamesdk.com/cron-task/internal/service/sql_cleaning"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"context"
	cron_task2 "github.com/cngamesdk/go-core/model/sql/cron_task"
)

// GetTaskFactory 获取任务工厂
func GetTaskFactory(taskType string) (resp TaskInterface) {
	switch taskType {
	case cron_task2.TaskTypeSqlCleaning:
		resp = &sql_cleaning.SqlCleaningService{}
		break
	}
	return
}

type TaskInterface interface {
	Init(req *cron_task.DimCronTaskConfigModel)
	PreEvent(ctx context.Context) (resp string, err error)
	Run(ctx context.Context) (err error)
	SuccessEvent(ctx context.Context) (err error)
	FailEvent(ctx context.Context) (err error)
	CompleteEvent(ctx context.Context) (err error)
}

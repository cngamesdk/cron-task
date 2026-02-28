package service

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"context"
	"github.com/cngamesdk/go-core/model/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type BaseService struct {
	IsRunning bool
	Config    *cron_task.DimCronTaskConfigModel
	TaskLog   *cron_task.OdsCronTaskLogModel
}

func (receiver *BaseService) Run(ctx context.Context) (err error) {
	if receiver.IsRunning {
		err = errors.New("正在执行中")
		return
	}
	//获取最新配置
	model := cron_task.NewDimCronTaskConfigModel()
	if takeErr := model.Take(ctx, "*", "id = ?", receiver.Config.Id); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取配置异常", zap.Error(takeErr))
		return
	}
	//已经下架
	if model.Status != sql.StatusNormal {
		err = errors.New("该任务已下架")
		global.Logger.ErrorCtx(ctx, "该任务已下架", zap.Any("data", model))
		return
	}
	//有变更
	if model.UpdatedAt != receiver.Config.UpdatedAt {
		receiver.Config = model
	}
	receiver.IsRunning = true
	return
}

func (receiver *BaseService) CompleteEvent(ctx context.Context) (err error) {
	receiver.IsRunning = false
	return receiver.TaskLog.Create(ctx)
}

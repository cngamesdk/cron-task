package sql_cleaning

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/internal/service/sql_cleaning"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

// AddFunc 增加SQL清洗任务
func AddFunc(myCron *cron.Cron, req *cron_task.DimCronTaskConfigModel) (cron.EntryID, error) {
	service := sql_cleaning.NewSqlCleaningService(req)
	return myCron.AddFunc(req.Spec, func() {
		requestId := cryptor.Md5String(fmt.Sprintf("%d%s", time.Now().UnixMilli(), random.RandString(5)))
		ctx := context.WithValue(context.Background(), global.Config.Common.CtxRequestIdKey, requestId)
		if service.IsRunning {
			global.Logger.WarnCtx(ctx, "正在执行中", zap.Any("data", req))
			return
		}
		global.Logger.InfoCtx(ctx, "开始执行", zap.Any("data", req.Name))
		startTime := time.Now()
		service.IsRunning = true
		defer func() {
			service.IsRunning = false
		}()
		runErr := service.Run(ctx)
		global.Logger.InfoCtx(ctx, "结束执行", zap.Any("time", time.Now().Sub(startTime).Seconds()))
		if runErr != nil {
			global.Logger.ErrorCtx(ctx, "执行异常", zap.Error(runErr))
		}
	})
}

package initialization

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"cngamesdk.com/cron-task/model/task"
	"context"
	"github.com/cngamesdk/go-core/model/sql"
	"github.com/duke-git/lancet/v2/random"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

var (
	tasksQueue = make(map[int64]int64)
)

// InitTasks 获取所有任务列表
func InitTasks(myCron *cron.Cron) {

	for {

		time.Sleep(time.Second * 10)

		model := cron_task.NewDimCronTaskConfigModel()
		tmpDb := model.Db().Table(model.TableName()).Select("*").Where("status = ?", sql.StatusNormal)
		var count int64
		if countErr := tmpDb.Count(&count).Error; countErr != nil {
			global.Logger.Error("获取总数异常", zap.Error(countErr))
			continue
		}
		if count <= 0 {
			global.Logger.Info("未获取到任务")
			continue
		}
		page := 1
		pageSize := 50
		totalPage := cast.ToInt(count) / pageSize
		if cast.ToInt(count)%pageSize != 0 {
			totalPage += 1
		}
		for page <= totalPage {
			var list []cron_task.DimCronTaskConfigModel
			if listErr := model.Db().
				Table(model.TableName()).
				Select("*").
				Where("status = ?", sql.StatusNormal).
				Limit(pageSize).
				Offset((page - 1) * pageSize).
				Order("id DESC").Find(&list).Error; listErr != nil {
				global.Logger.Error("获取列表异常", zap.Error(listErr))
				continue
			}
			for _, item := range list {

				if _, ok := tasksQueue[item.Id]; ok {
					continue
				}
				//新增
				tasksQueue[item.Id] = item.Id

				if item.Config == nil {
					item.Config = make(sql.CustomMapType)
				}
				adapter := task.GetTaskFactory(item.TaskType)
				if adapter == nil {
					global.Logger.Warn("未知任务类型", zap.Any("data", item))
					continue
				}
				adapter.Init(&item)

				entryId, addFunErr := myCron.AddFunc(item.Spec, func() {
					requestId, _ := random.UUIdV4()
					ctx := context.WithValue(context.Background(), global.Config.Common.CtxRequestIdKey, requestId)
					global.Logger.InfoCtx(ctx, "开始执行任务")
					if runErr := adapter.Run(ctx); runErr != nil {
						global.Logger.Error("执行任务异常", zap.Any("err", runErr))
					}
					global.Logger.InfoCtx(ctx, "结束执行任务")
				})

				if addFunErr != nil {
					global.Logger.Warn("加入定时异常", zap.Any("err", addFunErr))
					continue
				}
				global.Logger.Info("任务开启", zap.Any("实例ID", entryId), zap.Any("任务ID", item.Id))
			}
			page++
		}
	}
}

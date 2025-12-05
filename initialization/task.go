package initialization

import (
	"cngamesdk.com/cron-task/global"
	sql_cleaning2 "cngamesdk.com/cron-task/internal/logic/sql_cleaning"
	"cngamesdk.com/cron-task/model/sql/cron_task"
	"github.com/cngamesdk/go-core/model/sql"
	cron_task2 "github.com/cngamesdk/go-core/model/sql/cron_task"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// InitTasks 获取所有任务列表
func InitTasks(myCron *cron.Cron) (err error) {
	model := cron_task.NewDimCronTaskConfigModel()
	tmpDb := model.Db().Table(model.TableName()).Select("*").Where("status = ?", sql.StatusNormal)
	var count int64
	if countErr := tmpDb.Count(&count).Error; countErr != nil {
		err = countErr
		return
	}
	if count <= 0 {
		err = errors.New("未找到任务列表")
		return
	}
	page := 1
	pageSize := 50
	totalPage := cast.ToInt(count) / pageSize
	if cast.ToInt(count)%pageSize != 0 {
		totalPage++
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
			err = listErr
			return
		}
		for _, item := range list {
			var entryId cron.EntryID
			var addFunErr error
			if item.Config == nil {
				item.Config = make(sql.CustomMapType)
			}
			switch item.TaskType {
			case cron_task2.TaskTypeSqlCleaning: // SQL清洗任务
				entryId, addFunErr = sql_cleaning2.AddFunc(myCron, &item)
				break
			default:
				err = errors.New("未知任务类型" + item.TaskType)
				return
			}
			if addFunErr != nil {
				err = addFunErr
				return
			}
			global.Logger.Info("任务开始执行", zap.Any("任务ID", entryId))
		}
		page++
	}
	return
}

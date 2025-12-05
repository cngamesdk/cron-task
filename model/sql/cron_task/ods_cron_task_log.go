package cron_task

import (
	"cngamesdk.com/cron-task/global"
	"github.com/cngamesdk/go-core/model/sql/cron_task"
	"gorm.io/gorm"
)

type OdsCronTaskLogModel struct {
	cron_task.OdsCronTaskLogModel
}

func NewOdsCronTaskLogModel() *OdsCronTaskLogModel {
	model := &OdsCronTaskLogModel{}
	model.OdsCronTaskLogModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}

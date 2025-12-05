package cron_task

import (
	"cngamesdk.com/cron-task/global"
	"github.com/cngamesdk/go-core/model/sql/cron_task"
	"gorm.io/gorm"
)

type DimCronTaskConfigModel struct {
	cron_task.DimCronTaskConfigModel
}

func NewDimCronTaskConfigModel() *DimCronTaskConfigModel {
	model := &DimCronTaskConfigModel{}
	model.DimCronTaskConfigModel.Db = func() *gorm.DB {
		return global.MyDb
	}
	return model
}

package initialization

import (
	"cngamesdk.com/cron-task/global"
	"github.com/cngamesdk/go-core/model/sql/cron_task"
	"github.com/cngamesdk/go-core/model/sql/log"
)

// Migrate 迁移数据
func Migrate() (err error) {
	if migrateErr := global.MyDb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC").
		AutoMigrate(
			log.DwdGameRegLogModel{},
			log.DwdRootGameRegLogModel{},
			log.DwdRootGameBackRegLogModel{},

			cron_task.DimCronTaskConfigModel{},
			cron_task.OdsCronTaskLogModel{},
		); migrateErr != nil {
		err = migrateErr
	}
	return
}

package initialization

import (
	"cngamesdk.com/cron-task/config"
	"cngamesdk.com/cron-task/global"
	config3 "github.com/cngamesdk/go-core/config"
	log2 "github.com/cngamesdk/go-core/log"
)

func Init(config config.Config) {
	global.Logger = log2.MyLogger{
		CtxRequestIdKey: global.Config.Common.CtxRequestIdKey,
	}
	global.Logger.Logger = log2.NewFileZapLogger(config.Log)

	//初始化数据库
	db, dbErr := config3.OpenMysql(config.Mysql)

	if dbErr != nil {
		panic(any(dbErr))
	}
	global.MyDb = db

	if migrateErr := Migrate(); migrateErr != nil {
		panic(any(migrateErr))
	}

	//初始化REDIS
	myRedis, myRedisErr := config3.OpenRedis(config.Redis)
	if myRedisErr != nil {
		panic(any(myRedisErr))
	}
	global.MyRedis = myRedis
}

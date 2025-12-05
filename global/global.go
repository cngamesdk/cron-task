package global

import (
	"cngamesdk.com/cron-task/config"
	log2 "github.com/cngamesdk/go-core/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Logger  log2.MyLogger
	Config  config.Config
	MyDb    *gorm.DB
	MyRedis *redis.Client
)

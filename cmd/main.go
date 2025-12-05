package main

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/initialization"
	"cngamesdk.com/cron-task/logger"
	"flag"
	"github.com/robfig/cron/v3"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	var config string
	flag.StringVar(&config, "config", "", "-config=/your/config/path")
	flag.Parse()
	if config == "" {
		panic(any("配置不能为空"))
	}

	if initDataErr := initialization.InitConfigData(config); initDataErr != nil {
		panic(any(initDataErr))
	}

	defer global.Logger.Logger.Sync()

	initialization.Init(global.Config)

	cronLog := logger.CronLog{}
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(cronLog)))
	if err := initialization.InitTasks(c); err != nil {
		panic(any(err))
	}

	c.Start()
	select {}
}

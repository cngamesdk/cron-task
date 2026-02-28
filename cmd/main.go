package main

import (
	"cngamesdk.com/cron-task/global"
	"cngamesdk.com/cron-task/initialization"
	"cngamesdk.com/cron-task/logger"
	"flag"
	"github.com/cngamesdk/go-core/goroutine"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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

	//协程异步处理
	goroutine.CreateGoroutine(func() {
		cronLog := logger.CronLog{}
		c := cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.Recover(cronLog)))
		initialization.InitTasks(c)
		c.Start()
		println("任务开启.")
	}, func(any2 any) {
		global.Logger.Info("异步任务异常", zap.Any("err", any2))
	})
	println("执行中...")
	select {}
}

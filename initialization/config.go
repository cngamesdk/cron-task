package initialization

import (
	"cngamesdk.com/cron-task/global"
	"fmt"
	"github.com/spf13/viper"
)

func InitConfigData(config string) (err error) {
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if errRead := v.ReadInConfig(); errRead != nil {
		err = fmt.Errorf("Fatal error config file: %s", errRead)
		return
	}
	if errJson := v.Unmarshal(&global.Config); errJson != nil {
		err = fmt.Errorf("Fatal error unmarshal file: %s", err)
		return
	}

	return
}

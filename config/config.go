package config

import (
	"github.com/cngamesdk/go-core/config"
)

type Config struct {
	Common config.CommonConfig `mapstructure:"common" json:"common" yaml:"common"`
	Mysql  config.MySql        `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  config.Redis        `mapstructure:"redis" json:"redis" yaml:"redis"`
	Log    config.FileLog      `mapstructure:"log" json:"log" yaml:"log"`
}

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var GlobalConfig *viper.Viper

const ConfPath = "conf/config.yml"

func LoadConf(filepath string) {
	GlobalConfig = viper.New()
	GlobalConfig.SetConfigFile(filepath)
	err := GlobalConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}
}

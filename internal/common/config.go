package common

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/needon1997/theshop-svc/internal/common/config"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func LoadConfig(configPath string, devMode bool) {
	splits := strings.Split(configPath, ".")
	configType := splits[len(splits)-1]
	fileReader, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigType(configType)
	panicIfError(v.ReadConfig(fileReader))
	panicIfError(v.Unmarshal(&config.ServerConfig))
	if !devMode {
		port, err := GetFreePort()
		if err != nil {
			panic(err)
		}
		config.ServerConfig.Port = port
		config.ServerConfig.ConsulConfig.Id = fmt.Sprintf("%s:%v", config.ServerConfig.ConsulConfig.Id, port)
	}
	config.ServerConfig.ConsulConfig.Id = uuid.New().String()
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

package common

import (
	"github.com/needon1997/theshop-svc/internal/common/config"
	"go.uber.org/zap"
)

func NewLogger(devMode bool) {
	logPath := config.ServerConfig.LogConfig.LogPath
	var logger *zap.Logger
	var err error
	if devMode {
		logger, err = NewDevLogger(logPath)
	} else {
		logger, err = NewProductLogger(logPath)
	}
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func NewDevLogger(logPath string) (*zap.Logger, error) {
	if logPath == "" {
		return zap.NewDevelopment()
	}
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths = []string{logPath, "stdout"}
	return loggerConfig.Build()
}

func NewProductLogger(logPath string) (*zap.Logger, error) {
	if logPath == "" {
		return zap.NewProduction()
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{logPath, "stderr"}
	return loggerConfig.Build()

}

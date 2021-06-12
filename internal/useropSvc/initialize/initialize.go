package initialize

import (
	"github.com/needon1997/theshop-svc/internal/common"
	"github.com/needon1997/theshop-svc/internal/useropSvc/model"
	"go.uber.org/zap"
)

func Initialization() {
	ParseFlag()
	common.LoadConfig(*ConfigPath, *DevMode)
	common.NewLogger(*DevMode)
	err := common.RegisterSelfToConsul()
	if err != nil {
		zap.S().Errorw("Fail to register to consul", "error", err.Error)
	}
	common.PanicIfError(model.InitConnection())
}

func Finalize() {
	err := common.DeRegisterFromConsul()
	if err != nil {
		zap.S().Errorw("Fail to deregister from consul", "error", err.Error)
	}
	zap.S().Sync()
	zap.L().Sync()
}

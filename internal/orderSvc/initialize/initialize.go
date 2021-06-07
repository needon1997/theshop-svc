package initialize

import (
	"github.com/needon1997/theshop-svc/internal/common"
	"github.com/needon1997/theshop-svc/internal/common/grpc_client"
	"github.com/needon1997/theshop-svc/internal/inventorySvc/model"
	"github.com/needon1997/theshop-svc/internal/orderSvc/service"
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
	service.GoodsSvcConn, err = grpc_client.GetGoodsSvcConn()
	if err != nil {
		zap.S().Errorw("Fail to get goods svc connection", "error", err.Error)
	}
	service.InventorySvcConn, err = grpc_client.GetInventorySvcConn()
	if err != nil {
		zap.S().Errorw("Fail to get inventory goods svc connection", "error", err.Error)
	}
}

func Finalize() {
	err := common.DeRegisterFromConsul()
	if err != nil {
		zap.S().Errorw("Fail to deregister from consul", "error", err.Error)
	}
	zap.S().Sync()
	zap.L().Sync()
}

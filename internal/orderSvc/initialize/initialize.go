package initialize

import (
	"github.com/needon1997/theshop-svc/internal/common"
	"github.com/needon1997/theshop-svc/internal/common/grpc_client"
	"github.com/needon1997/theshop-svc/internal/orderSvc/global"
	"github.com/needon1997/theshop-svc/internal/orderSvc/model"
	"go.uber.org/zap"
	"io"
)

var traceCloser io.Closer

func Initialization() {
	ParseFlag()
	common.LoadConfig(*ConfigPath, *DevMode)
	common.NewLogger(*DevMode)
	traceCloser = common.InitJaeger()
	err := common.RegisterSelfToConsul()
	if err != nil {
		zap.S().Errorw("Fail to register to consul", "error", err.Error)
	}
	common.PanicIfError(model.InitConnection())
	global.GoodsSvcConn, err = grpc_client.GetGoodsSvcConn()
	if err != nil {
		zap.S().Errorw("Fail to get goods svc connection", "error", err.Error)
	}
	global.InventorySvcConn, err = grpc_client.GetInventorySvcConn()
	if err != nil {
		zap.S().Errorw("Fail to get inventory goods svc connection", "error", err.Error)
	}
	InitSqs()
}

func Finalize() {
	global.GoodsSvcConn.Close()
	global.InventorySvcConn.Close()
	traceCloser.Close()
	err := common.DeRegisterFromConsul()
	if err != nil {
		zap.S().Errorw("Fail to deregister from consul", "error", err.Error)
	}
	zap.S().Sync()
	zap.L().Sync()
}

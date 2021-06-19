package sqs

import (
	"context"
	"fmt"
	"github.com/needon1997/theshop-svc/internal/orderSvc/global"
	"github.com/needon1997/theshop-svc/internal/orderSvc/model"
	"github.com/needon1997/theshop-svc/internal/orderSvc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func panicIfErr(desc string, err error) {
	if err != nil {
		var wrapperErr error
		errDesc := fmt.Sprintf("[%s error]: %s", desc, err.Error())
		switch err {
		case gorm.ErrRecordNotFound:
			wrapperErr = status.Error(codes.NotFound, errDesc)
		default:
			wrapperErr = status.Error(codes.Internal, errDesc)
		}
		panic(wrapperErr)
	}
}
func HandleTimeOutOrder(in *proto.SellInfo) (err error) {
	tx := model.DB.Begin()
	defer func() {
		if r, ok := recover().(error); ok {
			tx.Rollback()
			zap.S().Infow("[HandleTimeOutOrder]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	order := model.OrderInfo{}
	panicIfErr("Find order", tx.Where("order_sn = ?", in.OrderSn).First(&order).Error)
	if order.Status == "approved" {
		return
	}
	order.Status = "timeout"
	inventoryClient := proto.NewInventoryClient(global.InventorySvcConn)
	_, err1 := inventoryClient.Reback(context.Background(), in)
	if err1 != nil {
		panic(err1)
	}
	panicIfErr("Order timeout", tx.Save(&order).Error)
	panicIfErr("Commit Transaction", tx.Commit().Error)
	return
}

package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/inventorySvc/model"
	"github.com/needon1997/theshop-svc/internal/inventorySvc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type InventoryService struct {
	proto.UnimplementedInventoryServer
}

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

func (InventoryService) SetInv(ctx context.Context, in *proto.GoodsInvInfo) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[SetInv]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	inventory := model.Inventory{}
	err = model.DB.Where("goods_id = ?", in.GoodsId).First(&inventory).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panicIfErr("find goods", err)
		}
		inventory.Stocks = uint(in.Num)
		inventory.GoodsID = uint(in.GoodsId)
		panicIfErr("Create Inventory", model.DB.Create(&inventory).Error)
		rsp = &empty.Empty{}
		err = nil
		return
	} else {
		inventory.Stocks = uint(in.Num)
		panicIfErr("Update inventory", model.DB.Save(&inventory).Error)
		rsp = &empty.Empty{}
		return
	}
}
func (InventoryService) InvDetail(ctx context.Context, in *proto.GoodsInvInfo) (rsp *proto.GoodsInvInfo, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[InvDetail]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	inventory := model.Inventory{}
	panicIfErr("Find Inventory", model.DB.Where("goods_id = ?", in.GoodsId).First(&inventory).Error)
	rsp = &proto.GoodsInvInfo{
		GoodsId: in.GoodsId,
		Num:     int32(inventory.Stocks),
	}
	return
}
func (InventoryService) Sell(ctx context.Context, in *proto.SellInfo) (rsp *empty.Empty, err error) {
	tx := model.DB.Begin()
	defer func() {
		if r, ok := recover().(error); ok {
			tx.Rollback()
			zap.S().Infow("[Sell]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	for _, goodsInfo := range in.GoodsInfo {
		inventory := model.Inventory{}
		panicIfErr("Get Stocks", tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inventory).Error)
		if inventory.Stocks >= uint(goodsInfo.Num) {
			inventory.Stocks -= uint(goodsInfo.Num)
			panicIfErr("Update Stocks", tx.Save(&inventory).Error)
		}
	}
	panicIfErr("commit change", tx.Commit().Error)
	rsp = &empty.Empty{}
	return
}
func (InventoryService) Reback(ctx context.Context, in *proto.SellInfo) (rsp *empty.Empty, err error) {
	tx := model.DB.Begin()
	defer func() {
		if r, ok := recover().(error); ok {
			tx.Rollback()
			zap.S().Infow("[Reback]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	for _, goodsInfo := range in.GoodsInfo {
		inventory := model.Inventory{}
		panicIfErr("Get Stocks", tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inventory).Error)
		inventory.Stocks += uint(goodsInfo.Num)
		panicIfErr("Update Stocks", tx.Save(&inventory).Error)
	}
	panicIfErr("commit change", tx.Commit().Error)
	rsp = &empty.Empty{}
	return
}

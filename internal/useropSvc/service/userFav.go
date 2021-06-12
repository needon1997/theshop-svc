package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/useropSvc/model"
	"github.com/needon1997/theshop-svc/internal/useropSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserFavService struct {
	proto.UnimplementedUserFavServer
}

func (UserFavService) GetFavList(ctx context.Context, in *proto.UserFavRequest) (rsp *proto.UserFavListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetFavList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	uf := make([]model.UserFav, 0)
	var db *gorm.DB = model.DB
	if in.UserId != 0 {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.GoodsId != 0 {
		db = db.Where("goods_id = ?", in.GoodsId)
	}
	panicIfErr("find userFav", db.Find(&uf).Error)
	rsp = &proto.UserFavListResponse{
		Total: int32(len(uf)),
	}
	for i := 0; i < int(rsp.Total); i++ {
		rsp.Data = append(rsp.Data, &proto.UserFavResponse{
			UserId:  int32(uf[i].UserId),
			GoodsId: int32(uf[i].GoodsId),
		})
	}
	return
}
func (UserFavService) AddUserFav(ctx context.Context, in *proto.UserFavRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[AddUserFav]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	uf := model.UserFav{
		UserId:  uint(in.UserId),
		GoodsId: uint(in.GoodsId),
	}
	panicIfErr("Create UserFav", model.DB.Create(&uf).Error)
	rsp = &empty.Empty{}
	return
}
func (UserFavService) DeleteUserFav(ctx context.Context, in *proto.UserFavRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteUserFav]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	uf := model.UserFav{}
	panicIfErr("Find userFav", model.DB.Where("user_id = ?", in.UserId).Where("goods_id = ?", in.GoodsId).First(&uf).Error)
	panicIfErr("Delete userFav", model.DB.Delete(&uf).Error)
	rsp = &empty.Empty{}
	return
}
func (UserFavService) GetUserFavDetail(ctx context.Context, in *proto.UserFavRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetUserFavDetail]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	uf := model.UserFav{}
	panicIfErr("Find userFav", model.DB.Where("user_id = ?", in.UserId).Where("goods_id = ?", in.GoodsId).First(&uf).Error)
	rsp = &empty.Empty{}
	return
}

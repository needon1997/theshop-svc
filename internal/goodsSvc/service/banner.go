package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (GoodsService) BannerList(ctx context.Context, in *empty.Empty) (rsp *proto.BannerListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[BannerList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	banners := make([]model.Banner, 0)
	panicIfErr(FIND_BANNER, model.DB.Find(&banners).Error)
	rsp = &proto.BannerListResponse{}
	rsp.Total = int32(len(banners))
	for i := 0; i < len(banners); i++ {
		rsp.Data = append(rsp.Data, &proto.BannerResponse{
			Id:    int32(banners[i].ID),
			Index: int32(banners[i].Index),
			Image: banners[i].Image,
			Url:   banners[i].Url,
		})
	}
	return
}
func (GoodsService) CreateBanner(ctx context.Context, in *proto.BannerRequest) (rsp *proto.BannerResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateBanner]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	banner := model.Banner{
		Model: gorm.Model{ID: uint(in.Id)},
		Image: in.Image,
		Url:   in.Url,
		Index: uint(in.Index),
	}
	panicIfErr(CREATE_BANNER, model.DB.Create(&banner).Error)
	rsp = &proto.BannerResponse{
		Id:    int32(banner.ID),
		Index: int32(banner.Index),
		Image: banner.Image,
		Url:   banner.Url,
	}
	return
}
func (GoodsService) DeleteBanner(ctx context.Context, in *proto.BannerRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteBanner]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_BANNER, model.DB.First(&model.Banner{}, in.Id).Error)
	panicIfErr(DELETE_BANNER, model.DB.Delete(&model.Banner{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) UpdateBanner(ctx context.Context, in *proto.BannerRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateBanner]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_BANNER, model.DB.First(&model.Banner{}, in.Id).Error)
	banner := model.Banner{
		Image: in.Image,
		Url:   in.Url,
		Index: uint(in.Index),
	}
	panicIfErr(UPDATE_BANNER, model.DB.Model(&model.Banner{Model: gorm.Model{ID: uint(in.Id)}}).Select("Image", "Url", "Index").Updates(&banner).Error)
	rsp = &empty.Empty{}
	return
}

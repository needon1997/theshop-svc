package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (GoodsService) BrandList(ctx context.Context, in *proto.BrandFilterRequest) (rsp *proto.BrandListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Errorw("[BrandList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	brands := make([]model.Brand, 0)
	panicIfErr(FIND_BRAND, model.DB.Find(&brands).Error)
	rsp = &proto.BrandListResponse{}
	rsp.Total = int32(len(brands))
	pn := 1
	psize := 10
	if in.Pages != 0 {
		pn = int(in.Pages)
	}
	if in.PagePerNums != 0 {
		psize = int(in.PagePerNums)
	}
	offset := (pn - 1) * psize
	for i := offset; i < offset+psize; i++ {
		rsp.Data = append(rsp.Data, &proto.BrandInfoResponse{
			Id:   int32(brands[i].ID),
			Name: brands[i].Name,
			Logo: brands[i].Logo,
		})
	}
	return
}
func (GoodsService) CreateBrand(ctx context.Context, in *proto.BrandRequest) (rsp *proto.BrandInfoResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Errorw("[CreateBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	brand := model.Brand{
		Model: gorm.Model{ID: uint(in.Id)},
		Name:  in.Name,
		Logo:  in.Logo,
	}
	panicIfErr(CREATE_BRAND, model.DB.Create(&brand).Error)
	rsp = &proto.BrandInfoResponse{
		Id:   int32(brand.ID),
		Name: brand.Name,
		Logo: brand.Logo,
	}
	return

}
func (GoodsService) DeleteBrand(ctx context.Context, in *proto.BrandRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Errorw("[DeleteBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_BRAND, model.DB.First(&model.Brand{}, in.Id).Error)
	panicIfErr(DELETE_BRAND, model.DB.Delete(&model.Brand{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) UpdateBrand(ctx context.Context, in *proto.BrandRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Errorw("[UpdateBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_BRAND, model.DB.First(&model.Brand{}, in.Id).Error)
	brand := model.Brand{
		Name: in.Name,
		Logo: in.Logo,
	}
	panicIfErr(UPDATE_BRAND, model.DB.Model(&model.Brand{Model: gorm.Model{ID: uint(in.Id)}}).Select("Name", "Logo").Updates(&brand).Error)
	rsp = &empty.Empty{}
	return
}

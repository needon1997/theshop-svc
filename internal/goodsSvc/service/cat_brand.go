package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (GoodsService) CategoryBrandList(ctx context.Context, in *proto.CategoryBrandFilterRequest) (rsp *proto.CategoryBrandListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CategoryBrandList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	catBrands := make([]model.GoodsBrandCategory, 0)
	panicIfErr(FIND_CAT_BRAND, model.DB.Preload("Category").Preload("Brand").Find(&catBrands).Error)
	rsp = &proto.CategoryBrandListResponse{}
	rsp.Total = int32(len(catBrands))
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
		rsp.Data = append(rsp.Data, &proto.CategoryBrandResponse{
			Id: int32(catBrands[i].ID),
			Brand: &proto.BrandInfoResponse{
				Id:   int32(catBrands[i].Brand.ID),
				Name: catBrands[i].Brand.Name,
				Logo: catBrands[i].Brand.Logo,
			},
			Category: &proto.CategoryInfoResponse{
				Id:             int32(catBrands[i].Category.ID),
				Name:           catBrands[i].Category.CategoryName,
				ParentCategory: int32(catBrands[i].Category.ParentCategoryID),
				Level:          int32(catBrands[i].Category.Level),
				IsTab:          catBrands[i].Category.IsShow,
			},
		})
	}
	return
}
func (GoodsService) GetCategoryBrandList(ctx context.Context, in *proto.CategoryInfoRequest) (rsp *proto.BrandListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetCategoryBrandList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	catBrands := make([]model.GoodsBrandCategory, 0)
	panicIfErr(FIND_BRAND, model.DB.Preload("Brand").Where("category_id = ?", in.Id).Find(&catBrands).Error)
	rsp = &proto.BrandListResponse{}
	rsp.Total = int32(len(catBrands))
	for i := 0; i < len(catBrands); i++ {
		rsp.Data = append(rsp.Data, &proto.BrandInfoResponse{
			Id:   int32(catBrands[i].Brand.ID),
			Name: catBrands[i].Brand.Name,
			Logo: catBrands[i].Brand.Logo,
		})
	}
	return
}
func (GoodsService) CreateCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (rsp *proto.CategoryBrandResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateCategoryBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	catBrand := model.GoodsBrandCategory{
		CategoryID: uint(in.CategoryId),
		BrandID:    uint(in.BrandId),
	}
	panicIfErr(FIND_CAT, model.DB.First(&model.Category{}, in.CategoryId).Error)
	panicIfErr(FIND_BRAND, model.DB.First(&model.Brand{}, in.BrandId).Error)
	panicIfErr(CREATE_CAT_BRAND, model.DB.Create(&catBrand).Error)
	rsp = &proto.CategoryBrandResponse{
		Id: int32(catBrand.ID),
	}
	return
}
func (GoodsService) DeleteCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteCategoryBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT_BRAND, model.DB.First(&model.GoodsBrandCategory{}, in.Id).Error)
	panicIfErr(DELETE_CAT_BRAND, model.DB.Delete(&model.GoodsBrandCategory{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) UpdateCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateCategoryBrand]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT_BRAND, model.DB.First(&model.GoodsBrandCategory{}, in.Id).Error)
	catBrand := model.GoodsBrandCategory{
		CategoryID: uint(in.CategoryId),
		BrandID:    uint(in.BrandId),
	}
	panicIfErr(UPDATE_CAT_BRAND, model.DB.Model(&model.GoodsBrandCategory{Model: gorm.Model{ID: uint(in.Id)}}).Select("CategoryID", "BrandID").Updates(&catBrand).Error)
	rsp = &empty.Empty{}
	return
}

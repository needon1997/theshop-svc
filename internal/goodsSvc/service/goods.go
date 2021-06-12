package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (GoodsService) GoodsList(ctx context.Context, in *proto.GoodsFilterRequest) (rsp *proto.GoodsListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GoodsList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	db := model.DB
	if in.TopCategory != 0 {
		cat := model.Category{Model: gorm.Model{ID: uint(in.TopCategory)}}
		panicIfErr(FIND_CAT, db.First(&cat).Error)
		ids := make([]uint, 0)
		switch cat.Level {
		case 1:
			panicIfErr(FIND_CAT, db.Raw("SELECT id FROM categories c1 WHERE c1.parent_category_id IN (SELECT c2.id FROM categories c2 WHERE c2.parent_category_id = ?)", cat.ID).Scan(&ids).Error)
		case 2:
			panicIfErr(FIND_CAT, db.Model(model.Category{}).Where("parent_category_id = ?", cat.Level).Select("id").Scan(&ids).Error)
		case 3:
			ids = append(ids, cat.ID)
		}
		db = db.Where("category_id IN ?", ids)
	}
	if in.KeyWords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", in.KeyWords))
	}
	if in.PriceMin != 0 {
		db = db.Where("shop_price >= ?", in.PriceMin)
	}
	if in.PriceMax != 0 {
		db = db.Where("shop_price <= ?", in.PriceMax)
	}
	c := model.Goods{IsHot: in.IsHot, IsNew: in.IsNew, IsShow: in.IsTab, BrandID: uint(in.Brand)}
	goods := make([]model.Goods, 0)
	panicIfErr(FIND_GOODS, db.Preload("Category").Preload("Brand").Where(&c).Find(&goods).Error)
	rsp = &proto.GoodsListResponse{}
	rsp.Total = int32(len(goods))
	psize := 10
	if in.PagePerNums != 0 {
		psize = int(in.PagePerNums)
	}
	pn := 1
	if in.Pages != 0 {
		pn = int(in.Pages)
	}
	offset := psize * (pn - 1)
	for i := offset; i < offset+psize; i++ {
		rsp.Data = append(rsp.Data, goodsModelToGoodsResponse(goods[i]))
	}
	return
}
func (GoodsService) BatchGetGoods(ctx context.Context, in *proto.BatchGoodsIdInfo) (rsp *proto.GoodsListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[BatchGetGoods]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	rsp = &proto.GoodsListResponse{}
	goods := make([]model.Goods, 0)
	panicIfErr(FIND_GOODS, model.DB.Find(&goods, in.Id).Error)
	rsp.Total = int32(len(goods))
	for i := 0; i < int(rsp.Total); i++ {
		rsp.Data = append(rsp.Data, goodsModelToGoodsResponse(goods[i]))
	}
	return
}
func (GoodsService) CreateGoods(ctx context.Context, in *proto.CreateGoodsInfo) (rsp *proto.GoodsInfoResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateGoods]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT, model.DB.First(&model.Category{}, in.CategoryId).Error)
	panicIfErr(FIND_BRAND, model.DB.First(&model.Brand{}, in.BrandId).Error)
	goods := createGoodsInfoToGoodsModel(in)
	panicIfErr(CREATE_GOODS, model.DB.Create(&goods).Error)
	rsp = goodsModelToGoodsResponse(goods)
	//TODO inventory service
	return
}
func (GoodsService) DeleteGoods(ctx context.Context, in *proto.DeleteGoodsInfo) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteGoods]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_GOODS, model.DB.First(&model.Goods{}, in.Id).Error)
	panicIfErr(DELETE_GOODS, model.DB.Delete(&model.Goods{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) UpdateGoods(ctx context.Context, in *proto.CreateGoodsInfo) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateGoods]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT, model.DB.First(&model.Category{}, in.CategoryId).Error)
	panicIfErr(FIND_BRAND, model.DB.First(&model.Brand{}, in.BrandId).Error)
	panicIfErr(FIND_GOODS, model.DB.First(&model.Goods{}, in.Id).Error)
	goods := createGoodsInfoToGoodsModel(in)
	panicIfErr(UPDATE_GOODS, model.DB.Model(&model.Goods{Model: gorm.Model{ID: uint(in.Id)}}).Select("*").Omit("ID", "CreatedAt").Updates(&goods).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) GetGoodsDetail(ctx context.Context, in *proto.GoodInfoRequest) (rsp *proto.GoodsInfoResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetGoodsDetail]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	good := model.Goods{}
	panicIfErr(FIND_GOODS, model.DB.Preload("Category").Preload("Brand").First(&good, in.Id).Error)
	panicIfErr(UPDATE_GOODS, model.DB.Model(&good).Update("click_num", gorm.Expr("click_num + ?", 1)).Error)
	good.ClickNum++
	rsp = goodsModelToGoodsResponse(good)
	return
}

package service

import (
	"fmt"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

const (
	JSON_PARSE       = "parse json"
	FIND_GOODS       = "find goods"
	UPDATE_GOODS     = "update goods"
	CREATE_GOODS     = "create goods"
	DELETE_GOODS     = "delete goods"
	FIND_BRAND       = "find brand"
	UPDATE_BRAND     = "update brand"
	CREATE_BRAND     = "create brand"
	DELETE_BRAND     = "delete brand"
	FIND_CAT         = "find category"
	UPDATE_CAT       = "update category"
	CREATE_CAT       = "create category"
	DELETE_CAT       = "delete category"
	FIND_BANNER      = "find banner"
	UPDATE_BANNER    = "update banner"
	CREATE_BANNER    = "create banner"
	DELETE_BANNER    = "delete banner"
	FIND_CAT_BRAND   = "find cat_banner"
	UPDATE_CAT_BRAND = "update cat_banner"
	CREATE_CAT_BRAND = "create cat_banner"
	DELETE_CAT_BRAND = "delete cat_banner"
)

type GoodsService struct {
	proto.UnimplementedGoodsServer
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

var imageDelimiter = "#1102011#"

func goodsModelToGoodsResponse(goods model.Goods) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              int32(goods.ID),
		CategoryId:      int32(goods.CategoryID),
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        int32(goods.ClickNum),
		SoldNum:         int32(goods.SoldNum),
		FavNum:          int32(goods.FavNum),
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.Brief,
		GoodsDesc:       goods.Description,
		ShipFree:        goods.ShipFree,
		Images:          strings.Split(goods.Images, imageDelimiter),
		DescImages:      strings.Split(goods.DescImages, imageDelimiter),
		GoodsFrontImage: goods.FrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.IsShow,
		AddTime:         goods.CreatedAt.Unix(),
		Category: &proto.CategoryBriefInfoResponse{
			Id:   int32(goods.Category.ID),
			Name: goods.Category.CategoryName,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   int32(goods.Brand.ID),
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

func createGoodsInfoToGoodsModel(in *proto.CreateGoodsInfo) model.Goods {
	return model.Goods{
		CategoryID:  uint(in.CategoryId),
		BrandID:     uint(in.BrandId),
		GoodsSn:     in.GoodsSn,
		Name:        in.Name,
		Stock:       uint(in.Stocks),
		MarketPrice: in.MarketPrice,
		ShopPrice:   in.ShopPrice,
		Brief:       in.GoodsBrief,
		Description: in.GoodsDesc,
		ShipFree:    in.ShipFree,
		Images:      strings.Join(in.Images, imageDelimiter),
		DescImages:  strings.Join(in.DescImages, imageDelimiter),
		FrontImage:  in.GoodsFrontImage,
		IsNew:       in.IsNew,
		IsHot:       in.IsHot,
		IsShow:      in.OnSale,
	}
}

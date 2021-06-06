package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/model"
	"github.com/needon1997/theshop-svc/internal/goodsSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (GoodsService) GetAllCategorysList(ctx context.Context, in *empty.Empty) (rsp *proto.CategoryListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetAllCategorysList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cats := make([]model.Category, 0)
	catsInfo := make([]*proto.CategoryInfoResponse, 0)
	panicIfErr(FIND_CAT, model.DB.Find(&cats).Error)
	l1 := make([]model.Category, 0)
	l2 := make([]model.Category, 0)
	l3 := make([]model.Category, 0)
	for i := 0; i < len(cats); i++ {
		catsInfo = append(catsInfo, &proto.CategoryInfoResponse{
			Id:             int32(cats[i].ID),
			Name:           cats[i].CategoryName,
			ParentCategory: int32(cats[i].ParentCategoryID),
			Level:          int32(cats[i].Level),
			IsTab:          cats[i].IsShow,
		})
		if cats[i].Level == 1 {
			l1 = append(l1, cats[i])
		}
		if cats[i].Level == 2 {
			l2 = append(l2, cats[i])
		}
		if cats[i].Level == 3 {
			l3 = append(l3, cats[i])
		}
	}
	catsMap1 := make([]map[string]interface{}, 0)
	for i := 0; i < len(l1); i++ {
		cm1 := make(map[string]interface{})
		bytes, err1 := json.Marshal(l1[i])
		panicIfErr(JSON_PARSE, err1)
		panicIfErr(JSON_PARSE, json.Unmarshal(bytes, &cm1))
		catsMap2 := make([]map[string]interface{}, 0)
		for j := 0; j < len(l2); j++ {
			if l2[j].ParentCategoryID == l1[i].ID {
				cm2 := make(map[string]interface{})
				bytes, err1 = json.Marshal(l2[j])
				panicIfErr(JSON_PARSE, err1)
				panicIfErr(JSON_PARSE, json.Unmarshal(bytes, &cm2))
				catsMap3 := make([]map[string]interface{}, 0)
				for k := 0; k < len(l3); k++ {
					if l3[k].ParentCategoryID == l2[j].ID {
						cm3 := make(map[string]interface{})
						bytes, err1 = json.Marshal(l3[k])
						panicIfErr(JSON_PARSE, err1)
						panicIfErr(JSON_PARSE, json.Unmarshal(bytes, &cm3))
						catsMap3 = append(catsMap3, cm3)
					}
				}
				cm2["sub_category"] = catsMap3
				catsMap2 = append(catsMap2, cm2)
			}
		}
		cm1["sub_category"] = catsMap2
		catsMap1 = append(catsMap1, cm1)
	}
	jsonString, err1 := json.Marshal(catsMap1)
	panicIfErr(JSON_PARSE, err1)
	rsp = &proto.CategoryListResponse{}
	rsp.Total = int32(len(cats))
	rsp.Data = catsInfo
	rsp.JsonData = string(jsonString)
	return
}
func (GoodsService) GetSubCategory(ctx context.Context, in *proto.CategoryListRequest) (rsp *proto.SubCategoryListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetSubCategory]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cat := model.Category{}
	panicIfErr(FIND_CAT, model.DB.First(&cat, in.Id).Error)
	sub_cat := make([]model.Category, 0)
	des := fmt.Sprintf("find sub_cat of cat:%v", cat.ID)
	panicIfErr(des, model.DB.Where("parent_category_id = ?", cat.ID).Find(&sub_cat).Error)
	rsp = &proto.SubCategoryListResponse{}
	rsp.Total = int32(len(sub_cat))
	rsp.Info = &proto.CategoryInfoResponse{
		Id:             int32(cat.ID),
		Name:           cat.CategoryName,
		ParentCategory: int32(cat.ParentCategoryID),
		Level:          int32(cat.Level),
		IsTab:          cat.IsShow,
	}
	for i := 0; i < len(sub_cat); i++ {
		rsp.SubCategorys = append(rsp.SubCategorys, &proto.CategoryInfoResponse{
			Id:             int32(sub_cat[i].ID),
			Name:           sub_cat[i].CategoryName,
			ParentCategory: int32(sub_cat[i].ParentCategoryID),
			Level:          int32(sub_cat[i].Level),
			IsTab:          cat.IsShow,
		})
	}
	return
}
func (GoodsService) CreateCategory(ctx context.Context, in *proto.CategoryInfoRequest) (rsp *proto.CategoryInfoResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateCategory]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	cat := model.Category{
		CategoryName:     in.Name,
		Level:            uint(in.Level),
		IsShow:           in.IsTab,
		ParentCategoryID: uint(in.ParentCategory),
	}
	panicIfErr(CREATE_CAT, model.DB.Create(&cat).Error)
	rsp = &proto.CategoryInfoResponse{
		Id:             int32(cat.ID),
		Name:           cat.CategoryName,
		ParentCategory: int32(cat.ParentCategoryID),
		Level:          int32(cat.Level),
		IsTab:          cat.IsShow,
	}
	return
}
func (GoodsService) DeleteCategory(ctx context.Context, in *proto.DeleteCategoryRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteCategory]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT, model.DB.First(&model.Category{}, in.Id).Error)
	panicIfErr(DELETE_CAT, model.DB.Delete(&model.Category{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (GoodsService) UpdateCategory(ctx context.Context, in *proto.CategoryInfoRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[UpdateCategory]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr(FIND_CAT, model.DB.First(&model.Category{}, in.Id).Error)
	cat := model.Category{
		CategoryName: in.Name,
		IsShow:       in.IsTab,
	}
	panicIfErr(UPDATE_CAT, model.DB.Model(&model.Category{Model: gorm.Model{ID: uint(in.Id)}}).Select("CategoryName", "IsShow").Updates(&cat).Error)
	rsp = &empty.Empty{}
	return
}

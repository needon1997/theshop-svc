package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/useropSvc/model"
	"github.com/needon1997/theshop-svc/internal/useropSvc/proto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddressService struct {
	proto.UnimplementedAddressServer
}

func (AddressService) GetAddressList(ctx context.Context, in *proto.AddressRequest) (rsp *proto.AddressListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[GetAddressList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	var db *gorm.DB = model.DB
	if in.UserId != 0 {
		db = model.DB.Where("user_id = ?", in.UserId)
	}
	addressList := make([]model.Address, 0)
	panicIfErr("find address", db.Find(&addressList).Error)
	rsp = &proto.AddressListResponse{Total: int32(len(addressList))}
	for i := 0; i < int(rsp.Total); i++ {
		rsp.Data = append(rsp.Data, &proto.AddressResponse{
			Id:           int32(addressList[i].ID),
			UserId:       addressList[i].UserId,
			Province:     addressList[i].Province,
			City:         addressList[i].City,
			Address:      addressList[i].Address,
			SignerName:   addressList[i].Name,
			SignerMobile: addressList[i].Mobile,
		})
	}
	return
}
func (AddressService) CreateAddress(ctx context.Context, in *proto.AddressRequest) (rsp *proto.AddressResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateAddress]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	address := model.Address{
		UserId:   in.UserId,
		Province: in.Province,
		City:     in.City,
		Address:  in.Address,
		Name:     in.SignerName,
		Mobile:   in.SignerMobile,
	}
	panicIfErr("Create Address", model.DB.Create(&address).Error)
	rsp = &proto.AddressResponse{
		Id:           int32(address.ID),
		UserId:       address.UserId,
		Province:     address.Province,
		City:         address.City,
		Address:      address.Address,
		SignerName:   address.Name,
		SignerMobile: address.Mobile,
	}
	return
}
func (AddressService) DeleteAddress(ctx context.Context, in *proto.AddressRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteAddress]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	panicIfErr("Find address", model.DB.First(&model.Address{}, in.Id).Error)
	panicIfErr("Delete Address", model.DB.Delete(&model.Address{}, in.Id).Error)
	rsp = &empty.Empty{}
	return
}
func (AddressService) UpdateAddress(ctx context.Context, in *proto.AddressRequest) (rsp *empty.Empty, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[DeleteAddress]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	address := &model.Address{Model: gorm.Model{ID: uint(in.Id)}}
	panicIfErr("Find address", model.DB.First(&address).Error)
	if in.Address != "" {
		address.Address = in.Address
	}
	if in.City != "" {
		address.City = in.City
	}
	if in.Province != "" {
		address.Province = in.Province
	}
	if in.SignerName != "" {
		address.Name = in.SignerName
	}
	if in.SignerMobile != "" {
		address.Mobile = in.SignerMobile
	}
	panicIfErr("Update Address", model.DB.Save(&address).Error)
	return
}

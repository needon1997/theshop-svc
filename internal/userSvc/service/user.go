package service

import (
	"context"
	"github.com/needon1997/theshop-svc/internal/userSvc/model"
	"github.com/needon1997/theshop-svc/internal/userSvc/proto"
	"github.com/needon1997/theshop-svc/internal/userSvc/utils"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"time"
)

const USER_NOT_EXIST = "User not exist"
const USER_ALREADY_EXIST = "User already exist"

type UserService struct {
	proto.UnimplementedUserSVCServer
}

func (this *UserService) GetUserList(ctx context.Context, in *proto.PageInfoRequest) (*proto.UserListInfoResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	rsp := &proto.UserListInfoResponse{}
	total, err := model.CountUser()
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	rsp.Total = uint64(total)
	var page uint32 = 1
	var pageSize uint32 = 20
	if in.PageNumber != 0 {
		page = in.PageNumber
	}
	if in.PageSize != 0 {
		pageSize = in.PageSize
	}
	s1 := opentracing.GlobalTracer().StartSpan("db_query", opentracing.ChildOf(parentSpan.Context()))
	users, err := model.GetUserByOffSetLimit(int(page*pageSize), int(pageSize))
	s1.Finish()
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	for _, user := range users {
		userInfo := &proto.UserInfoResponse{}
		mapUserInfo(user, userInfo)
		rsp.UserData = append(rsp.UserData, userInfo)
	}
	return rsp, nil
}
func (this *UserService) GetUserByEmail(ctx context.Context, in *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	email := in.Email
	user, err := model.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(codes.NotFound, USER_NOT_EXIST)
	}
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	userInfo := &proto.UserInfoResponse{}
	mapUserInfo(user, userInfo)
	return userInfo, nil
	return nil, nil
}
func (this *UserService) GetUserById(ctx context.Context, in *proto.IdRequest) (*proto.UserInfoResponse, error) {
	id := in.Id
	user, err := model.GetUserById(id)
	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(codes.NotFound, USER_NOT_EXIST)
	}
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	userInfo := &proto.UserInfoResponse{}
	mapUserInfo(user, userInfo)
	return userInfo, nil
}
func (this *UserService) CreateUser(ctx context.Context, in *proto.CreateUserInfoRequest) (*proto.UserInfoResponse, error) {
	user, err := model.GetUserByEmail(in.Email)
	if err != gorm.ErrRecordNotFound {
		if err != nil {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		} else {
			return nil, status.Error(codes.AlreadyExists, USER_ALREADY_EXIST)
		}
	}
	user = model.User{}
	user.Email = in.Email
	user.NickName = in.NickName
	user.Password = string(utils.EncryptPassword(in.Password))
	user, err = model.SaveUser(user)
	if err != nil {
		log.Println(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	userInfo := &proto.UserInfoResponse{}
	mapUserInfo(user, userInfo)
	return userInfo, nil
}
func (this *UserService) UpdateUser(ctx context.Context, in *proto.UpdateUserInfoRequest) (*proto.Empty, error) {
	user, err := model.GetUserById(in.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			log.Println(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	user.NickName = in.NickName
	user.Address = in.Address
	user.Birthday.Scan(time.Unix(int64(in.Birthday), 0))
	user.Gender = model.Gender(in.Gender)
	err = model.UpdateUser(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (this *UserService) ComparePassword(ctx context.Context, in *proto.ComparePasswordRequest) (*proto.BoolResponse, error) {
	err := bcrypt.CompareHashAndPassword([]byte(in.EncryptPwd), []byte(in.Password))
	if err != nil {
		return &proto.BoolResponse{Result: false}, nil
	}
	return &proto.BoolResponse{Result: true}, nil
}

func mapUserInfo(user model.User, userInfo *proto.UserInfoResponse) {
	userInfo.Id = uint64(user.ID)
	userInfo.NickName = user.NickName
	userInfo.Password = user.Password
	userInfo.Email = user.Email
	userInfo.Gender = string(user.Gender)
	userInfo.Address = user.Address
	if user.Birthday.Valid {
		userInfo.Birthday = uint64(user.Birthday.Time.Unix())
	}
	userInfo.Role = int32(user.Role)
}

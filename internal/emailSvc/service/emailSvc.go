package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/needon1997/theshop-svc/internal/emailSvc/model"
	"github.com/needon1997/theshop-svc/internal/emailSvc/proto"
	"github.com/needon1997/theshop-svc/internal/emailSvc/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EmailService struct {
	proto.UnimplementedEmailSvcServer
}

const VERIFICATION_SUBJECT = "Account Verification Code"
const VERIFICATION_TEMPLATE = "Your account verification code is %s"
const KEY_PREFIX = "email:"

func (this *EmailService) SendVerificationCode(ctx context.Context, in *proto.ReceiverInfoRequest) (*empty.Empty, error) {
	_, err := model.GetCodeByEmail(KEY_PREFIX + in.Email)
	if err != nil {
		if err == redis.Nil {
			code := utils.GenerateVerificationCode()
			err1 := model.SaveEmailCode(KEY_PREFIX+in.Email, code)
			if err1 != nil {
				zap.S().Errorw("Save Email Code Error", "error", err1.Error())
				return nil, status.Error(codes.Internal, "internal error")
			}
			go utils.SendVerificationEmail(in.Email, VERIFICATION_SUBJECT, fmt.Sprintf(VERIFICATION_TEMPLATE, code))
			return &empty.Empty{}, nil
		} else {
			zap.S().Errorw("Get Email Code Error", "error", err.Error())
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	code := utils.GenerateVerificationCode()
	err = model.SaveEmailCode(KEY_PREFIX+in.Email, code)
	if err != nil {
		zap.S().Errorw("Save Email Code Error", "error", err.Error())
		return nil, status.Error(codes.Internal, "internal error")
	}
	go utils.SendVerificationEmail(in.Email, VERIFICATION_SUBJECT, fmt.Sprintf(VERIFICATION_TEMPLATE, code))
	return &empty.Empty{}, nil
}
func (this *EmailService) VerifyVerificationCode(ctx context.Context, in *proto.VerifyCodeRequest) (*proto.VerifyResponse, error) {
	code, err := model.GetCodeByEmail(KEY_PREFIX + in.Email)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Error(codes.NotFound, "email record not found")
		}
		zap.S().Errorw("Get Email Code Error", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &proto.VerifyResponse{Match: code == in.Code}, nil
}
func (this *EmailService) SendMarketingEmail(context.Context, *proto.MarketingInfoRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMarketingEmail not implemented")
}
func (this *EmailService) SendTransactionalEmail(context.Context, *proto.TransactionalInfoRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransactionalEmail not implemented")
}

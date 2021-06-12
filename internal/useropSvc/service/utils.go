package service

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

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

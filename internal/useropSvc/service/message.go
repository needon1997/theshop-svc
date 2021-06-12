package service

import (
	"context"
	"github.com/needon1997/theshop-svc/internal/useropSvc/model"
	"github.com/needon1997/theshop-svc/internal/useropSvc/proto"
	"go.uber.org/zap"
)

type MessageService struct {
	proto.UnimplementedMessageServer
}

func (MessageService) MessageList(c context.Context, in *proto.MessageRequest) (rsp *proto.MessageListResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[MessageList]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	db := model.DB
	if in.UserId != 0 {
		db.Where("user_id = ?", in.UserId)
	}
	messageList := make([]model.UserMessage, 0)
	panicIfErr("Find userMessage", db.Find(&messageList).Error)
	rsp = &proto.MessageListResponse{Total: int32(len(messageList))}
	for i := 0; i < int(rsp.Total); i++ {
		rsp.Data = append(rsp.Data, &proto.MessageResponse{
			Id:          int32(messageList[i].ID),
			UserId:      messageList[i].UserId,
			MessageType: messageList[i].MessageType,
			Subject:     messageList[i].Subject,
			Message:     messageList[i].Content,
			File:        messageList[i].File,
		})
	}
	return
}
func (MessageService) CreateMessage(c context.Context, in *proto.MessageRequest) (rsp *proto.MessageResponse, err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			zap.S().Infow("[CreateMessage]", "[ERROR]:", r.Error())
			err = r
		}
	}()
	message := model.UserMessage{
		MessageType: in.MessageType,
		Subject:     in.Subject,
		Content:     in.Message,
		File:        in.File,
		UserId:      in.UserId,
	}
	panicIfErr("Create Message", model.DB.Create(&message).Error)
	rsp = &proto.MessageResponse{
		Id:          int32(message.ID),
		UserId:      message.UserId,
		MessageType: message.MessageType,
		Subject:     message.Subject,
		Message:     message.Content,
		File:        message.File,
	}
	return
}

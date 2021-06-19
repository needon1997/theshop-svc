package sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/needon1997/theshop-svc/internal/orderSvc/global"
	"github.com/needon1997/theshop-svc/internal/orderSvc/proto"
	"go.uber.org/zap"
)

func SendTimeoutOrderMessage(msg *proto.SellInfo) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = global.SqsSvc.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody:  aws.String(string(msgByte)),
		QueueUrl:     global.TimeoutOrderQueueName,
		DelaySeconds: 15 * 60, //15 minutes
	})
	return err
}

func PollAndHandleTimeoutOrderMessage() {
	for {
		msgResult, err := global.SqsSvc.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            global.TimeoutOrderQueueName,
			MaxNumberOfMessages: 1,
			VisibilityTimeout:   20,
		})
		if err != nil {
			zap.S().Errorw("[ReceiveMessage]", "Error", err)
			continue
		}
		if msgResult.Messages == nil {
			continue
		}
		handleString := msgResult.Messages[0].ReceiptHandle
		body := msgResult.Messages[0].Body
		err = HandleTimeoutOrder(*body)
		if err != nil {
			continue
		}
		_, err = global.SqsSvc.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
			QueueUrl:      global.TimeoutOrderQueueName,
			ReceiptHandle: handleString,
		})
		if err != nil {
			zap.S().Errorw("[AckMessage]", "Error", err)
			continue
		}
	}
}

func HandleTimeoutOrder(orderInfoString string) error {
	orderInfo := &proto.SellInfo{}
	err := json.Unmarshal([]byte(orderInfoString), orderInfo)
	if err != nil {
		zap.S().Errorw("[Parse SQS Message]", "Error", err)
		return err
	}
	err = HandleTimeOutOrder(orderInfo)
	return err
}

package initialize

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/needon1997/theshop-svc/internal/orderSvc/global"
	"go.uber.org/zap"
)

func InitSqs() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	svc := sqs.NewFromConfig(cfg)
	timeoutQueueUrlResult, err := svc.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: timeoutOrderQueueName,
	})
	if err != nil {
		zap.S().Errorw("[Initialize timeout order Sqs]", "ERROR", err)
	}
	global.TimeoutOrderQueueName = timeoutQueueUrlResult.QueueUrl
	global.SqsSvc = svc
}

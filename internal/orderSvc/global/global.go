package global

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"google.golang.org/grpc"
)

var TimeoutOrderQueueName *string
var SqsSvc *sqs.Client
var GoodsSvcConn *grpc.ClientConn
var InventorySvcConn *grpc.ClientConn

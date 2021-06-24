package grpc_client

import (
	"errors"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/needon1997/theshop-svc/internal/common/config"
	_ "github.com/needon1997/theshop-svc/internal/common/grpc_resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

const INTERNAL_ERROR = "server internal error"
const CONSUL_LB_TEMPLATE = "consul://%s/%s"

var opts = []grpc_retry.CallOption{
	grpc_retry.WithMax(3),
	grpc_retry.WithPerRetryTimeout(3 * time.Second),
	grpc_retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded, codes.Unknown),
}

func GetUserSvcConn() (*grpc.ClientConn, error) {
	zap.S().Debug("Get connect gRPC user service server")
	url := fmt.Sprintf(CONSUL_LB_TEMPLATE, config.ServerConfig.ServiceConfig.UserServiceName, "")
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		zap.S().Errorf("[GetUserSvcClient]  [fail to connect with service provider]   ERROR: %s", err.Error())
		return nil, errors.New(INTERNAL_ERROR)
	}
	return conn, nil
}
func GetEmailSvcConn() (*grpc.ClientConn, error) {
	zap.S().Debug("Get connect gRPC email service server")
	url := fmt.Sprintf(CONSUL_LB_TEMPLATE, config.ServerConfig.ServiceConfig.EmailServiceName, "")
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		zap.S().Errorf("[GetEmailSvcClient]  [fail to connect with service provider]   ERROR: %s", err.Error())
		return nil, errors.New(INTERNAL_ERROR)
	}
	return conn, nil
}
func GetGoodsSvcConn() (*grpc.ClientConn, error) {
	zap.S().Debug("Get connect gRPC goods service server")
	url := fmt.Sprintf(CONSUL_LB_TEMPLATE, config.ServerConfig.ServiceConfig.GoodsServiceName, "")
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		zap.S().Errorf("[GetUserSvcClient]  [fail to connect with service provider]   ERROR: %s", err.Error())
		return nil, errors.New(INTERNAL_ERROR)
	}
	return conn, nil
}
func GetInventorySvcConn() (*grpc.ClientConn, error) {
	zap.S().Debug("Get connect gRPC goods service server")
	url := fmt.Sprintf(CONSUL_LB_TEMPLATE, config.ServerConfig.ServiceConfig.InventoryServiceName, "")
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		zap.S().Errorf("[GetUserSvcClient]  [fail to connect with service provider]   ERROR: %s", err.Error())
		return nil, errors.New(INTERNAL_ERROR)
	}
	return conn, nil
}
func GetOrderSvcConn() (*grpc.ClientConn, error) {
	zap.S().Debug("Get connect gRPC goods service server")
	url := fmt.Sprintf(CONSUL_LB_TEMPLATE, config.ServerConfig.ServiceConfig.OrderServiceName, "")
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		zap.S().Errorf("[GetUserSvcClient]  [fail to connect with service provider]   ERROR: %s", err.Error())
		return nil, errors.New(INTERNAL_ERROR)
	}
	return conn, nil
}

package main

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/needon1997/theshop-svc/internal/common"
	"github.com/needon1997/theshop-svc/internal/common/config"
	"github.com/needon1997/theshop-svc/internal/useropSvc/initialize"
	"github.com/needon1997/theshop-svc/internal/useropSvc/proto"
	"github.com/needon1997/theshop-svc/internal/useropSvc/service"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.Initialization()
	defer initialize.Finalize()
	if config.ServerConfig.Host == "" {
		panic("host not defined")
	}
	if config.ServerConfig.Port == 0 {
		panic("port not defined")
	}
	var opt []grpc.ServerOption
	//opt = append(opt, (grpc.UnaryInterceptor(AuthInterceptor)))
	opt = append(opt, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))
	server := grpc.NewServer(opt...)
	proto.RegisterMessageServer(server, service.MessageService{})
	proto.RegisterAddressServer(server, service.AddressService{})
	proto.RegisterUserFavServer(server, service.UserFavService{})
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		zap.S().Infof("server listen at %s:%v\n", config.ServerConfig.Host, config.ServerConfig.Port)
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", config.ServerConfig.Host, config.ServerConfig.Port))
		if err != nil {
			panic("port listen failure")
		}
		common.PanicIfError(server.Serve(lis))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Infow("shut down sever")
}

package main

import (
	"fmt"
	"github.com/needon1997/theshop-svc/internal/common"
	config2 "github.com/needon1997/theshop-svc/internal/common/config"
	"github.com/needon1997/theshop-svc/internal/emailSvc/initialize"
	"github.com/needon1997/theshop-svc/internal/emailSvc/proto"
	"github.com/needon1997/theshop-svc/internal/emailSvc/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.Initialization()
	defer initialize.Finalize()
	if config2.ServerConfig.Host == "" {
		panic("host not defined")
	}
	if config2.ServerConfig.Port == 0 {
		panic("port not defined")
	}
	var opt []grpc.ServerOption
	//opt = append(opt, (grpc.UnaryInterceptor(AuthInterceptor)))
	server := grpc.NewServer(opt...)
	proto.RegisterEmailSvcServer(server, &service.EmailService{})
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		zap.S().Infof("server listen at %s:%v\n", config2.ServerConfig.Host, config2.ServerConfig.Port)
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", config2.ServerConfig.Host, config2.ServerConfig.Port))
		if err != nil {
			panic("port listen failure")
		}
		common.PanicIfError(server.Serve(lis))
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shut down sever")
}

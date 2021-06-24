package grpc_resolver

import (
	"context"
	"fmt"
	"github.com/needon1997/theshop-svc/internal/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"sort"
	"time"
)

type ConsulBuilder struct {
}

func (b *ConsulBuilder) Scheme() string {
	return "consul"
}
func (b *ConsulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	svcChan := make(chan []common.Service)
	go getServices(ctx, svcChan, target)
	go updateServices(ctx, svcChan, cc, target.Authority)
	return &resolvr{cancelFunc: cancel}, nil
}

const SERVICES_FILTER_URI = "http://%s/v1/agent/services?filter=Service==\"%s\""

func getServices(ctx context.Context, svcChan chan<- []common.Service, target resolver.Target) {
	svcBus := make(chan []common.Service)
	quit := make(chan bool)
	go func() {
		for {
			svcs, err := common.GetServicesByNameTags(target.Authority, target.Endpoint)
			if err != nil {
				zap.S().Errorw("Refresh Service Error", "Error: ", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}
			select {
			case <-quit:
				return
			case svcBus <- svcs:
				continue
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			quit <- true
			return
		case svcs := <-svcBus:
			svcChan <- svcs
		}
	}
}

func updateServices(ctx context.Context, svcChan <-chan []common.Service, cc resolver.ClientConn, serviceName string) {
	for {
		select {
		case svcs := <-svcChan:
			addrs := make([]resolver.Address, 0)
			for i := 0; i < len(svcs); i++ {
				addrs = append(addrs, resolver.Address{Addr: fmt.Sprintf("%s:%v", svcs[i].Address, svcs[i].Port)})
			}
			sort.Sort(byAddressString(addrs))
			zap.S().Infow("update service address", "Service: ", serviceName)
			cc.UpdateState(resolver.State{Addresses: addrs})
			time.Sleep(20 * time.Second)
		case <-ctx.Done():
			return
		}
	}
}

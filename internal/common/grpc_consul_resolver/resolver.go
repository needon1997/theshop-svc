package grpc_consul_resolver

import (
	"context"
	"fmt"
	"github.com/needon1997/theshop-svc/internal/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"sort"
	"time"
)

func init() {
	resolver.Register(&ConsulBuilder{})
}

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

type byAddressString []resolver.Address

func (p byAddressString) Len() int           { return len(p) }
func (p byAddressString) Less(i, j int) bool { return p[i].Addr < p[j].Addr }
func (p byAddressString) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

type resolvr struct {
	cancelFunc context.CancelFunc
}

// ResolveNow will be skipped due unnecessary in this case
func (r *resolvr) ResolveNow(resolver.ResolveNowOptions) {}

// Close closes the resolver.
func (r *resolvr) Close() {
	r.cancelFunc()
}

//func main() {
//	grpc.Dial("consul://localhost:8500/", grpc.WithInsecure())
//	for {
//
//	}
//}

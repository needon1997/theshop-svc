package common

import (
	"fmt"
	config1 "github.com/needon1997/theshop-svc/internal/common/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

func InitJaeger() io.Closer {
	cfg := &config.Configuration{
		ServiceName: config1.ServerConfig.TraceConfig.Name,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%v", config1.ServerConfig.TraceConfig.Host, config1.ServerConfig.TraceConfig.Port),
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}

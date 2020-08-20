package xjaeger

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var (
	_jaegerAgentEndpoint = "172.16.70.1:6832"
	//_jaegerRecorderEndpoint = "http://localhost:14268/api/traces"
)

func BootJaegerTracer(localServiceName, hostPort string) (opentracing.Tracer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		ServiceName: localServiceName,
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			//CollectorEndpoint: _jaegerRecorderEndpoint,
			LocalAgentHostPort:  _jaegerAgentEndpoint,
			BufferFlushInterval: 5 * time.Second,
		},
	}

	tracer, _, err := cfg.NewTracer(
		config.Logger(jaegerlog.StdLogger),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "BootJaegerTracer")
	}

	return tracer, nil
}

func GetTraceIdFromSpanContext(spanCtx opentracing.SpanContext) string {
	sc, ok := spanCtx.(jaeger.SpanContext)
	if ok {
		return sc.TraceID().String()
	}

	return ""
}

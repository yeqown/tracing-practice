package xjaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	_jaegerRecorderEndpoint = "https://localhost:14268/api/traces"
)

func BootJaegerTracer(localServiceName, hostPort string) (opentracing.Tracer, error) {
	cfg := &config.Configuration{
		ServiceName: localServiceName,
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: _jaegerRecorderEndpoint,
		},
	}

	tracer, _, err := cfg.NewTracer()
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

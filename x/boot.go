package x

import (
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	xzipkin "github.com/yeqown/opentracing-practice/x/x-zipkin"
)

func BootTracerWrapper(localServiceName string, hostPort string) error {
	tracer, err := xzipkin.BootZipkinTracer(localServiceName, hostPort)
	if err != nil {
		return errors.Wrap(err, "BootTracerWrapper.BootZipkinTracer")
	}
	opentracing.SetGlobalTracer(tracer)

	return nil
}

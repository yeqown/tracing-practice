package x

//
//import (
//	"github.com/opentracing/opentracing-go"
//	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
//	"github.com/openzipkin/zipkin-go"
//	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
//	"github.com/pkg/errors"
//)
//
//var (
//	_zipkinRecorderEndpoint = "http://localhost:9411/api/v2/spans"
//)
//
//func BootZipkinTracer(localServiceName, hostPort string) (opentracing.Tracer, error) {
//	reporter := zipkinhttp.NewReporter(
//		_zipkinRecorderEndpoint,
//		// zipkinhttp.Logger(log.New(os.Stdout, "zipkin", log.Lshortfile)),
//		// zipkinhttp.BatchInterval(100*time.Millisecond),
//	)
//	localEndpoint, err := zipkin.NewEndpoint(localServiceName, hostPort)
//	if err != nil {
//		return nil, errors.Wrap(err, "zipkin.NewEndpoint")
//	}
//	nativeTracer, err := zipkin.NewTracer(
//		reporter,
//		zipkin.WithTraceID128Bit(false), // TODO: diff between 128 and 64bit
//		zipkin.WithLocalEndpoint(localEndpoint),
//		zipkin.WithSharedSpans(true),
//		// TODO: more options
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	tracer := zipkinot.Wrap(nativeTracer)
//	return tracer, nil
//}
//
//func GetTraceIdFromSpanContext(spanCtx opentracing.SpanContext) string {
//	sc, ok := spanCtx.(zipkinot.SpanContext)
//	if ok {
//		return sc.TraceID.String()
//	}
//
//	return ""
//}

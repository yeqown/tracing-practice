package main

import (
	"log"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	endpointUrl = "http://localhost:9411/api/v2/spans"
	name        = "http"
)

func bootTracer() (*zipkin.Tracer, error) {
	localEndpoint, err := zipkin.NewEndpoint("service-a", "127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}

	reporter := zipkinhttp.NewReporter(endpointUrl)
	nativeTracer, err := zipkin.NewTracer(reporter,
		zipkin.WithTraceID128Bit(false),
		zipkin.WithSharedSpans(false),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)

	return nativeTracer, nil
}

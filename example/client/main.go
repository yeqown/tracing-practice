package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yeqown/opentracing-practice/x"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	// _ "github.com/openzipkin/zipkin-go/middleware/grpc"
)

var (
	serverAddr = "http://127.0.0.1:8080"
)

func bootstrap() {
	if err := x.BootTracerWrapper(
		"client", "127.0.0.1:50000"); err != nil {
		log.Fatal(err)
	}
}

func main2() {
	bootstrap()

	tracer := opentracing.GlobalTracer()
	if tracer == nil {
		panic("tracer not set")
	}

	// generate span
	_, sp := x.DeriveFromContext(context.Background())
	defer sp.Finish()

	// HTTP Client
	r := resty.New().R()
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	err := tracer.Inject(sp.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		panic(err)
	}

	// do request
	resp, err := r.Get(serverAddr + "/trace")
	if err != nil {
		log.Fatal(err)
	}

	// read response
	fmt.Printf("%s\n", resp.Body())

	// FIXME:
	time.Sleep(5 * time.Second)
}

//
//func main() {
//	bootstrap()
//
//	tracer := opentracing.GlobalTracer()
//	if tracer == nil {
//		panic("tracer not set")
//	}
//
//	// TODO: ignore zipkinTracer
//	client, err := zipkinhttp.NewClient(zipkinTracer, zipkinhttp.ClientTrace(true))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	req, _ := http.NewRequest(http.MethodGet, serverAddr+"/trace", nil)
//	resp, err := client.DoWithAppSpan(req, "client-request")
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = resp
//
//	println("Response: ", resp.StatusCode)
//
//	time.Sleep(5 * time.Second)
//}

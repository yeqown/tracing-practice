package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yeqown/opentracing-practice/x"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

func main() {
	bootstrap()

	// generate span
	_, sp := x.StartSpanFromContext(context.Background())
	defer func() {
		sp.Finish()
		time.Sleep(200 * time.Millisecond) // wait reporter to report
	}()

	url := serverAddr + "/trace"

	ext.SpanKindRPCClient.Set(sp)
	ext.HTTPUrl.Set(sp, url)
	ext.HTTPMethod.Set(sp, "GET")

	// HTTP Client
	r := resty.New().R()
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	if err := opentracing.GlobalTracer().
		Inject(sp.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		panic(err)
	}

	// do request
	resp, err := r.Get(url)
	if err != nil {
		panic(err)
	}

	// read response
	fmt.Printf("%s\n", resp.Body())
}

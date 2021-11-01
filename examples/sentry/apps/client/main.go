package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/go-resty/resty/v2"
)

var (
	serverAddr = "http://127.0.0.1:8080"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://1c2d1ae347944688ae7593a33e40c0f2@sentry.example.com/33",
		ServerName:  "client-demo",
		Environment: "dev",
		Release:     "v1.0.0",
		SampleRate:  1.0,
	})
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	sp := sentry.StartSpan(
		context.Background(),
		"request",
		sentry.TransactionName("from-client"),
	)
	defer sp.Finish()
	sp.Sampled = sentry.SampledTrue
	println(sp.ToSentryTrace())

	// HTTP Client
	r := resty.New().
		R().
		SetHeader("sentry-trace", sp.ToSentryTrace())

	// do request
	resp, err := r.Get(serverAddr + "/trace")
	if err != nil {
		panic(err)
	}

	// read response
	fmt.Printf("%s\n", resp.Body())
}

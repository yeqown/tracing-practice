package ccontext

import (
	"testing"
	"time"

	"github.com/yeqown/opentracing-practice/x"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger-client-go"
)

func Test_customContext_Value(t *testing.T) {
	ctx := newContext()
	ctx.Field1 = "1"

	err := x.BootTracerWrapper("127.0.0.1", "6381")
	require.Nil(t, err)

	sp1 := opentracing.StartSpan("test")
	ctx.Context = opentracing.ContextWithSpan(ctx.Context, sp1)
	sp2 := opentracing.SpanFromContext(ctx)
	assert.Equal(t, sp1, sp2)

	traceId1 := sp1.Context().(jaeger.SpanContext).TraceID()
	traceId2 := sp2.Context().(jaeger.SpanContext).TraceID()
	assert.Equal(t, traceId1, traceId2)
	t.Log(traceId1, traceId2)
}

func Test_customContext_WithTimeout(t *testing.T) {
	ctx := newContext()
	ctx = ctx.WithTimeout(5 * time.Second)
	done := make(chan struct{})

	go func() {
		defer func() {
			done <- struct{}{}
		}()

		tt := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-ctx.Done():
				t.Log("timeout")
				return
			case <-tt.C:
				t.Log("done")
				return
			}
		}
	}()

	<-done
}

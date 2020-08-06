# opentracing-practice
opentracing practice in golang micro server (gRPC + HTTP). I'm not using any standalone tools to trace, 
but part of them(zipkin/jaeger) implementation of opentracing appoint. 

## Practice trace chain
```sh
                                +-- process internal trace2
                                |
                     +---> process internal trace1
                     |
                     |                 +---> server-b trace(gRPC)
entry(HTTP) ---> server-a trace--gRPC--|
                                       +---> server-c trace(gRPC)
                                                   |
                                                   +----> process internal trace3
```

## Get started

* [practice with zipkin](./docs/zipkin-get-started.md)
* [practice with jaeger](./docs/jaeger-get-started.md)

## Conclusion [WIP] 🚀

### how to using opentracing?

first of all, you must boot an opentracing tracer, then register it into global (this is not necessary, but practical).

```go
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
```

* ***cross process***(cross servers)
    * `HTTP` client side.
    
        if you are doing a REST request and want to trace it, you may need to start span here.  
        
        ```go
        package main
        
        func main() {
            // ...
        
            // generate span
            _, sp := x.StartSpanFromContext(context.Background())
            defer sp.Finish()
            
            // set up request related info(URI, Method)
            ext.SpanKindRPCClient.Set(sp)
            ext.HTTPUrl.Set(sp, req.URI)
            ext.HTTPMethod.Set(sp, req.Method)
            
            // inject into request headere
            carrier := opentracing.HTTPHeadersCarrier(req.Header)
            err := opentracing.GlobalTracer().Inject(sp.Context(), opentracing.HTTPHeaders, carrier)
            
            // ...
        }
        ```   

    * `HTTP` middleware for server side.
    
        middleware need to do: 
        * parse parent span if HTTP client carried to you.
        * create a root span to pass by.
        
        ```go
        type getTraceID func(spCtx opentracing.SpanContext) string
        
        // get trace info from header, if not then create an new one
        func Opentracing(getTraceIdFromSpanContext getTraceID) gin.HandlerFunc {
            return func(c *gin.Context) {
                // prepare work ...
                carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
                clientSpCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
                if err == nil && clientSpCtx != nil {
                    sp = tracer.StartSpan(
                        c.Request.RequestURI,
                        opentracing.ChildOf(clientSpCtx),
                    )
                } else {
                    sp = tracer.StartSpan(c.Request.RequestURI)
                }
                defer sp.Finish()
                
                // do some work
                // ...
                
                ctx = opentracing.ContextWithSpan(c.Request.Context(), sp)
                c.Set(_traceContextKey, ctx)
                traceId := getTraceIdFromSpanContext(sp.Context())
                c.Header("X-Trace-Id", traceId)
                
                // continue process request
                c.Next()
                
                // do some work 
                // ...
            }
        }
        ```

    * `gRPC` interceptor (client and server).
        
        [TODO](#)
    
* ***internal process***(in one server)
    * derive the span

    ```go
    func StartSpanFromContext(ctx context.Context) (context.Context, opentracing.Span) {
        opName := WhoCalling()
        // StartSpanFromContext will create a span if ctx do not contains trace data. 
        sp, ctx := opentracing.StartSpanFromContext(ctx, opName)
        return ctx, sp
    }
    ```

### how opentracing works?

* [https://opentracing.io/docs/overview/](https://opentracing.io/docs/overview/)
* [https://sematext.com/blog/how-opentracing-works/](https://sematext.com/blog/how-opentracing-works/)

### What should I do if want to implement opentracing appoint? [TODO]

## References

* https://zhuanlan.zhihu.com/p/79419529
* https://opentracing.io/docs/getting-started/
* https://www.jaegertracing.io/docs/1.18/getting-started/
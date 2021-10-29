## OpenTracing support for gRPC in Go

>
> this part of code comes from `"github.com/yeqown/tracing-practice/x/opentracing/grpc-interceptor"/go/otgrpc`, was been modified for personal demand.
>

The `opentracingrpc` package makes it easy to add OpenTracing support to gRPC-based
systems in Go.

### Client-side usage example

Wherever you call `grpc.Dial`:

```go
// You must have some sort of OpenTracing Tracer instance on hand.
var tracer opentracing.Tracer = ...
...

// Set up a connection to the server peer.
conn, err := grpc.Dial(
    address,
    ... // other options
    grpc.WithUnaryInterceptor(
        opentracingrpc.OpenTracingClientInterceptor(tracer)),
    grpc.WithStreamInterceptor(
        opentracingrpc.OpenTracingStreamClientInterceptor(tracer)))

// All future RPC activity involving `conn` will be automatically traced.
```

### Server-side usage example

Wherever you call `grpc.NewServer`:

```go
// You must have some sort of OpenTracing Tracer instance on hand.
var tracer opentracing.Tracer = ...
...

// Initialize the gRPC server.
s := grpc.NewServer(
    ... // other options
    grpc.UnaryInterceptor(
        opentracingrpc.OpenTracingServerInterceptor(tracer)),
    grpc.StreamInterceptor(
        opentracingrpc.OpenTracingStreamServerInterceptor(tracer)))

// All future RPC activity involving `s` will be automatically traced.
```


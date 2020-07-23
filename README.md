# opentracing-practice
opentracing practice in golang micro server (gRPC + HTTP)

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
                                                   +----> process internal trace3 (todo)
```

## Get started

1.Install `zipkin` and run.
```sh
docker run -d -p 9411:9411 \
--name zipkin \
docker.io/openzipkin/zipkin
```

2.Start server.

```sh
cd path/to/opentracing-practice

# server-c
go run example/server-c/main.go

# server-b
go run example/server-b/main.go

# server-a
go run example/server-b/main.go

# http-server
go run example/http/main.go
```

OR 

```
cd path/to/opentracing-practice

make run
```

3.Client do request.

```shell script
curl http://127.0.0.1:8080/trace
```

OR.
visit [http://127.0.0.1:8080/trace](http://127.0.0.1:8080/trace) in your browser.

4.Get traceId.

    from `conosole log` or `client response header`

## Result shots

![shot](./static/shot1.jpg)

## References

* https://zhuanlan.zhihu.com/p/79419529
* https://opentracing.io/docs/getting-started/
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

![zipkin-shot1](docs/shot1.png)

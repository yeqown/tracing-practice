echo:
	echo "opentracing-practice"

gen:
	protoc -I ./protos/  --go_out=plugins=grpc:./protogen  ./protos/ping_a.proto
	protoc -I ./protos/  --go_out=plugins=grpc:./protogen  ./protos/ping_b.proto
	protoc -I ./protos/  --go_out=plugins=grpc:./protogen  ./protos/ping_c.proto

run:
	go run ./example/server-c &
	go run ./example/server-b &
	go run ./example/server-a &

	go run ./example/http &

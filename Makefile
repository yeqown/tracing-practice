
protogen:
	protoc -I ./protos/  --go_out=plugins=grpc:./protogen  ./protos/ping.proto

run:
	go run ./example/server-c &
	go run ./example/server-b &
	go run ./example/server-a &

	go run ./example/http &

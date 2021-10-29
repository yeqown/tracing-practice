echo:
	echo "opentracing-practice"

run-opentracing:
	go run ./examples/opentracing/apps/c &
	go run ./examples/opentracing/apps/b &
	go run ./examples/opentracing/apps/a &

	go run ./examples/opentracing/apps/http &

echo:
	echo "opentracing-practice"

run-opentracing:
	go run ./opentracing/apps/c &
	go run ./opentracing/apps/b &
	go run ./opentracing/apps/a &

	go run ./opentracing/apps/http &

run-sentry:
	go run ./sentry/apps/c/main.go &
	go run ./sentry/apps/b/main.go &
	go run ./sentry/apps/a/main.go &

	go run ./sentry/apps/http/main.go &
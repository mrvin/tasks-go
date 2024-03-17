build:
	go build -o ../../bin/identidock -ldflags '-w -s'
run:
	SERVER_HTTP_ADDR="identidock:8888" \
	REDIS_CACHE_ADDR="redis:6379" \
	DNMONSTER_ADDR="dnmonster:8080" \
	go run main.go

.PHONY: build run

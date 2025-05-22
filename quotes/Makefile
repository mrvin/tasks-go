run:
	env `cat configs/quotes.env | xargs` ./bin/quotes-server
lint:
	golangci-lint run ./...

.PHONY: run lint


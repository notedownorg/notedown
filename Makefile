all: format test

format:
	gofmt -w .

test:
	go test -v ./...

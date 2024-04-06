all: format test dirty

dirty:
	git diff --exit-code

format:
	gofmt -w .

test:
	go test -v ./...

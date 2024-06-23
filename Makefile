all: format mod test dirty

dirty:
	git diff --exit-code

mod:
	go mod tidy

format:
	gofmt -w .

test:
	go test -v ./...

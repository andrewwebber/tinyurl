all: build test

build:
	mkdir -p ./bin || true
	go build -o ./bin ./...
test:
	go test -cover -race -v ./...
db:
	./hack/couchbase.sh


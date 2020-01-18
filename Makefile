test:
	go test -cover -race -v ./...
db:
	./hack/couchbase.sh


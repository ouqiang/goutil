export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: test
test:
	go test -race -v ./...


.PHONY: bench
bench:
	go test -race -v --bench=. ./...

.PHONY: test-cover
test-conver:
	go test -race -v ./... -cover

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:. httpclient/*.proto
export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: test
test:
	go test -race -v ./...


.PHONY: bench
bench:
	go test -race -run none -v --bench=. ./...

.PHONY: coverage
coverage:
	go test -run none -v -cover ./...

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:. httpclient/*.proto
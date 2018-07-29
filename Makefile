
export CGO_ENABLED=1

.PHONY: test
test:
	go test -race ./...


.PHONY: bench
bench:
	go test -race --bench=. ./...

.PHONY: test-cover
test-conver:
	go test -race ./... -cover
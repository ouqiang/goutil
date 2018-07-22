
export CGO_ENABLED=1

.PHONY: test
test:
	go test -race ./... -v

.PHONY: bench
bench:
	go test -race --bench=. ./... -v

.PHONY: test-cover
test-conver:
	go test -race ./... -cover
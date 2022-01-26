test:
	go test ./pkg/... ./cmd/...

vet:
	go vet ./pkg/... ./cmd/...

check: test vet

build:
	go build -o target/check_diff ./cmd/...
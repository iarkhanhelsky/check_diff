test:
	go test ./pkg/... ./cmd/... ./integration_tests/...

vet:
	go vet ./pkg/... ./cmd/...

check: test vet

build:
	go build -o target/check_diff ./cmd/...
test:
	go test -v -covermode atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

vet:
	go vet ./pkg/... ./cmd/...

check: test vet

build:
	go build -o target/check_diff .
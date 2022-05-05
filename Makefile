MOCKGEN=go run github.com/golang/mock/mockgen
mocks/mock_shell.go: pkg/core/shell.go
	${MOCKGEN} -destination=mocks/mock_shell.go \
               -package=mocks \
               github.com/iarkhanhelsky/check_diff/pkg/core Shell

mockgen: mocks/mock_shell.go

test: mockgen
	go test -v -covermode atomic -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out

test_all: mockgen
	go test -v -covermode atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

vet:
	go vet ./pkg/... ./cmd/...

check: test vet

build:
	go build -o target/check_diff .
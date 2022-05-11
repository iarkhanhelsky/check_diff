MOCKGEN=go run github.com/golang/mock/mockgen

mocks/%_mock.go: $(wildcard $(dir %)/*.go)
	${MOCKGEN} \
 		-destination=$@ \
        -package=mocks \
        $(shell echo github.com/iarkhanhelsky/check_diff/$(subst mocks/,,$(dir $@)) | ruby -e 'puts gets[0...-2]' ) \
        $(shell echo $(subst _mock,,$(basename $(notdir $@))) |  ruby -e 'puts gets.capitalize')

mock_clean:
	mkdir -p mocks
	find mocks -name '*_mock.go' -exec rm {} \;

mockgen: mock_clean \
		 mocks/pkg/core/checker_mock.go \
		 mocks/pkg/unpack/unpacker_mock.go \
		 mocks/pkg/shell/shell_mock.go \
         mocks/pkg/tools/registry_mock.go \

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
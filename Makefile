.PHONY: build test test-race lint mock clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)
BIN := my-raindrop-io-mcp-server

build:
	go build -ldflags '$(LDFLAGS)' -o $(BIN) .

install:
	GOBIN=~/.local/bin go install -ldflags '$(LDFLAGS)' .

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run

mock:
	go run go.uber.org/mock/mockgen@latest \
		-source=internal/domain/repository/bookmark.go \
		-destination=internal/domain/repository/mock/bookmark_mock.go \
		-package=mock
	go run go.uber.org/mock/mockgen@latest \
		-source=internal/domain/repository/collection.go \
		-destination=internal/domain/repository/mock/collection_mock.go \
		-package=mock

clean:
	rm -f $(BIN)

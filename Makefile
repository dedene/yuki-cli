BINARY := bin/yuki
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w \
	-X github.com/dedene/yuki-cli/internal/cmd.Version=$(VERSION) \
	-X github.com/dedene/yuki-cli/internal/cmd.Commit=$(COMMIT) \
	-X github.com/dedene/yuki-cli/internal/cmd.Date=$(DATE)

CGO_ENABLED ?= $(shell [ "$$(uname)" = "Darwin" ] && echo 1 || echo 0)

.PHONY: build run test lint install clean fmt fmt-check ci

build:
	@mkdir -p $(dir $(BINARY))
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/yuki

run: build
	./$(BINARY) $(ARGS)

test:
	go test ./...

lint:
	GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint golangci-lint run

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/yuki

clean:
	@if [ -e "$(BINARY)" ]; then trash "$(BINARY)"; fi

fmt:
	gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then goimports -w .; fi
	@if command -v gofumpt >/dev/null 2>&1; then gofumpt -w .; fi

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Run 'make fmt'" && exit 1)
	@if command -v goimports >/dev/null 2>&1; then test -z "$$(goimports -l .)" || (echo "Run 'make fmt'" && exit 1); fi
	@if command -v gofumpt >/dev/null 2>&1; then test -z "$$(gofumpt -l .)" || (echo "Run 'make fmt'" && exit 1); fi

ci: fmt-check lint test

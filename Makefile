BINARY  := jviz
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags '-X main.version=$(VERSION)'

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: build build-all test lint check clean help

## build: Build for current platform
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY) .

## build-all: Build for all platforms
build-all:
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d/ -f1); \
		arch=$$(echo $$platform | cut -d/ -f2); \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		out="dist/$(BINARY)_$${os}_$${arch}$${ext}"; \
		echo "Building $$out…"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o "$$out" . || exit 1; \
	done

## test: Run tests
test:
	go test ./...

## lint: Run linter
lint:
	golangci-lint run

## check: Run tests and lint
check: test lint

## package: Create zip archives for all platforms
package: build-all
	@cd dist && for f in $(BINARY)_*; do \
		base=$$(basename $$f); \
		zip "$${base%.*}.zip" "$$base" ../README.md 2>/dev/null || zip "$${base}.zip" "$$base" ../README.md; \
	done
	@echo "Packages created in dist/"

## clean: Remove build artifacts
clean:
	rm -rf dist/

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //'

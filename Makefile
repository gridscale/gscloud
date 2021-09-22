VERSION := $(shell if ! git describe --tags 2>/dev/null; then \
	grep -Po '(?<=^VERSION=)v.*$$' $$PWD/RELEASE.txt; \
fi; \
)

GIT_COMMIT := $(shell if ! git rev-list -1 HEAD 2>/dev/null; then \
	grep -Po '(?<=^GIT_COMMIT=)\w*$$' $$PWD/RELEASE.txt; \
fi; \
)

default: build

build:
	go build -ldflags "-X github.com/gridscale/gscloud/cmd.GitCommit=$(GIT_COMMIT) -X github.com/gridscale/gscloud/cmd.Version=$(VERSION)"

test: build
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

lint: build
	golint ./...

clean:
	go clean

.PHONY: build test lint clean

VERSION := $(shell if ! git describe --tags 2>/dev/null; then \
	grep '^VERSION=' $$PWD/RELEASE.txt|cut -d'=' -f2-; \
fi; \
)

GIT_COMMIT := $(shell if ! git rev-list -1 HEAD 2>/dev/null; then \
	grep '^GIT_COMMIT=' $$PWD/RELEASE.txt|cut -d'=' -f2-; \
fi; \
)

default: build

build:
	go build -ldflags "-X github.com/gridscale/gscloud/cmd.GitCommit=$(GIT_COMMIT) -X github.com/gridscale/gscloud/cmd.Version=$(VERSION)"

build-debug:
	go build \
		-gcflags="all=-N -l" \
		-ldflags "-X github.com/gridscale/gscloud/cmd.GitCommit=$(GIT_COMMIT) -X github.com/gridscale/gscloud/cmd.Version=$(VERSION)"

test: build
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

lint: build
	staticcheck ./...

clean:
	go clean

.PHONY: build test lint clean

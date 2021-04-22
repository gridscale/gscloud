VERSION=$$(git describe --tags)
GIT_COMMIT=$$(git rev-list -1 HEAD)

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

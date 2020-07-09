PLATFORMS=windows linux darwin
ARCHES=amd64
BUILDDIR=build
VERSION=$$(cat VERSION)
GIT_COMMIT=$$(git rev-list -1 HEAD)
EXECUTABLE_NAME=gscloud_$(VERSION)

default: build

buildall: clean release zip

build:
	go build -ldflags "-X github.com/gridscale/gscloud/cmd.GitCommit=$(GIT_COMMIT) -X github.com/gridscale/gscloud/cmd.Version=$(VERSION)"

test: build
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

release:
	$(foreach platform,$(PLATFORMS), \
		$(foreach arch,$(ARCHES), \
			mkdir -p $(BUILDDIR); \
			GOOS=$(platform) GOARCH=$(arch) go build \
				-ldflags "-X github.com/gridscale/gscloud/cmd.GitCommit=$(GIT_COMMIT) -X github.com/gridscale/gscloud/cmd.Version=$(VERSION)" \
				-o $(BUILDDIR)/$(EXECUTABLE_NAME)_$(platform)_$(arch);))
	@echo "Renaming Windows file"
	@if [ -f $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES) ]; then mv $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES) $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe; fi

zip:
	$(foreach file,$(wildcard $(BUILDDIR)/*), \
		zip -j $(file).zip $(file);)
	@if [ -f $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe.zip ]; then mv $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe.zip $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).zip; fi

clean:
	go clean
	rm -f $(BUILDDIR)/gscloud_*

.PHONY: buildall build release clean zip

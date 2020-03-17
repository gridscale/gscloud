PLATFORMS=windows linux darwin
ARCHES=amd64
BUILDDIR=build
VERSION=$$(cat VERSION)
EXECUTABLE_NAME=gscloud_$(VERSION)

buildall: clean build zip

build:
	$(foreach platform,$(PLATFORMS), \
	    $(foreach arch,$(ARCHES), \
	        mkdir -p $(BUILDDIR); GOOS=$(platform) GOARCH=$(arch) go build -o $(BUILDDIR)/$(EXECUTABLE_NAME)_$(platform)_$(arch);))
	@echo "Renaming Windows file"
	@if [ -f $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES) ]; then mv $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES) $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe; fi

zip:
	$(foreach file,$(wildcard $(BUILDDIR)/*),\
		zip -j $(file).zip $(file);)
	@if [ -f $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe.zip ]; then mv $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).exe.zip $(BUILDDIR)/$(EXECUTABLE_NAME)_windows_$(ARCHES).zip; fi

clean:
	$(foreach platform,$(PLATFORMS), \
            $(foreach arch,$(ARCHES), \
                rm -f $(BUILDDIR)/gscloud_*;))

.PHONY: buildall build clean zip

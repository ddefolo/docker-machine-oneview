build-clean:
	rm -Rf $(PREFIX)/bin/*

extension = $(patsubst windows,.exe,$(filter windows,$(1)))

# Cross builder helper
define gocross
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build \
		-o $(PREFIX)/bin/docker-$(patsubst cmd/%.go,%,$3)_$(1)-$(2)$(call extension,$(GOOS)) \
		-a $(VERBOSE_GO) -tags "static_build netgo $(BUILDTAGS)" -installsuffix netgo \
		-ldflags "$(GO_LDFLAGS) -extldflags -static" $(GO_GCFLAGS) $(3);
endef


# get oneview golang vendor path
# we can make GO_LIB_ONEVIEW=/src/path:github.com/HewlettPackard/oneview-golang
# to setup gopath from src
GO_LIB_ONEVIEW ?= github.com/HewlettPackard/oneview-golang

define go-get
	if [ "$$(echo $(1) | cut -c1-1)" = "/" ]; then \
	  lib_source=$$(echo $(1)| awk -F':' '{print $$1}'); \
	  lib_path=$$(echo $(1)| awk -F':' '{print $$2}'); \
	  go_path=$$(echo $(GOPATH) | awk -F':' '{print $$1}'); \
	  cd $$go_path/src; \
	  mkdir -p $$lib_path; \
	  cd $$lib_path; \
	  source_count=$$(($$(find $$lib_source/* | wc -l)+1)); \
	  echo "Files in target dir --> $$(find . |wc -l)"; \
	  echo "Files in source dir --> $$source_count"; \
	  echo "Copy File in -- $$lib_source/ --> $$(pwd)"; \
	  cp -R $$lib_source/* .; \
	  target_count=$$(find . | wc -l); \
	  echo "New Files in target --> $$target_count"; \
	  if [ ! $$source_count = $$target_count ]; then \
	      echo "ERROR in setting up source dir, file count mis-match"; \
	      exit 1; \
	  fi; \
	else \
	  echo "installing $(1)"; \
	  go get -u $(1); \
	fi;
endef

go-install-oneview:
	@$(call go-get,$(GO_LIB_ONEVIEW))

# XXX building with -a fails in debug (with -N -l) ????

# Independent targets for every bin
$(PREFIX)/bin/docker-%: ./cmd/%.go $(shell find . -type f -name '*.go')
	$(GO) build -o $@$(call extension,$(GOOS)) $(VERBOSE_GO) -tags "$(BUILDTAGS)" -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $<

# Cross-compilation targets
build-x-%: ./cmd/%.go $(shell find . -type f -name '*.go')
	$(foreach GOARCH,$(TARGET_ARCH),$(foreach GOOS,$(TARGET_OS),$(call gocross,$(GOOS),$(GOARCH),$<)))

# Build all plugins
build-plugins: $(patsubst ./cmd/%.go,$(PREFIX)/bin/docker-%,$(filter-out %_test.go, $(wildcard ./cmd/machine-driver-*.go)))

# Overall cross-build
build-x: $(patsubst ./cmd/%.go,build-x-%,$(filter-out %_test.go, $(wildcard ./cmd/*.go)))

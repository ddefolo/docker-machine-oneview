# setup any dependencies for Build

GO_PACKAGES := github.com/docker/docker/pkg/term github.com/docker/machine github.com/stretchr/testify/assert
GO_PACKAGES := $(GO_PACKAGES) golang.org/x/crypto/ssh
# GOPATH := $(HOME)/go
# PATH := $(PATH):$(GOPATH)/bin:/usr/local/go/bin
# GO15VENDOREXPERIMENT := 1

# Cross builder helper
define godeps-get
	godep get $(1);
endef

define godeps-save
	godep save $(1);
endef

define GOVENDORPATH
$(shell godep path)
endef

define godeps-clean
	echo 'Clean Package $(1)';
	[ -d $(GOPATH)/src/$(1) ] && \
		( cd $(GOPATH)/src/$(1); \
			_PKG_CLEAN=`git rev-parse --show-toplevel`; \
			[ -d $$_PKG_CLEAN ] && rm -rf $$_PKG_CLEAN; ) || \
		echo "Skipting clean for $(1)";
endef

vendor-clean:
		@rm -rf $(PREFIX)/vendor/*
		@echo cleaning up in $(PREFIX)/vendor/*

# for fresh setup so we can do godep save -r
godeps-clean: vendor-clean
		@echo "Removing all dependent packages from $(GOPATH)"
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-clean,$(GOPCKG)))
		rm -rf $(GOPATH)/src/github.com/$(GH_USER)/$(GH_REPO)

# setup a fresh GOPATH directory with what would be needed to build
godeps-init: godeps-clean
		@echo "Pulling required packages into $(GOPATH)"
		mkdir -p $(GOPATH)/src/github.com/$(GH_USER)
		ln -s $(PREFIX) $(GOPATH)/src/github.com/$(GH_USER)/$(GH_REPO)
		@echo "Get dependent packages"
		# TODO: remove this call when we OSRB
		$(call godeps-get,github.com/HewlettPackard/oneview-golang)
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-get,$(GOPCKG)))

godeps-save:
		# TODO: add oneview-golang after OSRB $(call godeps-save, $(GO_PACKAGES) github.com/HewlettPackard/oneview-golang)
		$(call godeps-save, $(GO_PACKAGES))

# setup the vendor folder with required packages that have been committed
godeps-vendor:
		echo "Placing packages into $(GOVENDORPATH)"
		cp -R $(GOPATH)/src/github.com/HewlettPackard/oneview-golang/ $(PREFIX)/Godeps/_workspace/src/github.com/HewlettPackard/oneview-golang
		GOPATH=$(GOVENDORPATH) godep restore

godeps: godeps-init godeps-save
		echo "All done! If all looks good vendor it with : make godeps-vendor"

godep: godeps

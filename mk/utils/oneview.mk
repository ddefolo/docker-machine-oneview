define ONEVIEW_ROOT
$(PREFIX)/../oneview-golang
endef
define ONEVIEW_ROOT_CD
$(shell cd $(PREFIX)/../oneview-golang;pwd)
endef

golink-oneview-golang:
	[ ! -d "$(ONEVIEW_ROOT)" ] && \
		echo "ERROR: could not find $(ONEVIEW_ROOT)  \
		Try cloning the repo to $(ONEVIEW_ROOT)" && \
		exit 1; \
	[ ! -L $(PREFIX)/Godeps/_workspace/src/github.com/HewlettPackard/oneview-golang ] && \
		ln -s $(ONEVIEW_ROOT_CD) $(PREFIX)/Godeps/_workspace/src/github.com/HewlettPackard/oneview-golang; \
	exit 0;

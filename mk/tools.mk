# build projects in tools folder
# current tools include :
# oneview-machine - containerized version of this project

# copy binaries to tools/oneview-machine/bin
define tools-cpbin
    mkdir -p "$(PREFIX)/tools/oneview-machine/bin"; \
    for file in $$(find "$(PREFIX)/bin" -name 'docker-machine-driver-oneview_linux-amd64'); do \
        echo "copy -> $$file"; \
        cp "$$file" "$(PREFIX)/tools/oneview-machine/bin/"; \
    done;
endef

# build the oneview-machine container
define tools-oneview-machine
   set -x -v ; \
   cd "$(PREFIX)/tools/oneview-machine" && \
   make build
endef

clean-tools:
	@echo "cleaning binaries for tools container"
	rm -rf "$(PREFIX)/tools/oneview-machine/bin/docker-machine-driver*"
	rm -rf "$(PREFIX)/bin"

# lets first make some binaries
# then we'll copy them over and use them to build the oneview-machine
# container with latest bins in them

build-tools: clean-tools build
	@echo "building tools"
	@$(call tools-cpbin)
	@$(call tools-oneview-machine)


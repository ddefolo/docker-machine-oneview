PREFIX ?= $(shell pwd)

#
# simple function to alias with : eal $(make alias)
#
define alias-machine
    function oneview-machine { \
        envfile=$$(mktemp); \
        docker_env=$$(bash -c $(PREFIX)/env-opts.sh $$envfile); \
        args=$$@; \
        entrypoint="/usr/local/bin/docker-machine"; \
        if [ -z "$$args" ]; then entrypoint="bash"; fi; \
        eval docker run -it --rm \
                       "$$docker_env" \
                       --entrypoint $$entrypoint \
                       oneview-machine $$debug_opt $$args; \
        rm -f $$envfile; \
    }
endef

default: build

build:
	docker build --build-arg https_proxy=$(HTTPS_PROXY) \
	             --build-arg http_proxy=$(HTTP_PROXY) \
	             --build-arg HTTP_PROXY=$(HTTP_PROXY) \
	             --build-arg HTTPS_PROXY=$(HTTPS_PROXY) \
	             --tag oneview-machine .


#alias:
#	@echo alias oneview-machine=\'docker run -it --rm oneview-machine bash\'
#
alias:
	@echo '$(alias-machine)'
	@echo ''
	@echo '# source this script by running the command: eval $$(make alias)'
	@echo '# then you can run the command: oneview-machine'
	@echo ''

PREFIX ?= $(shell pwd)

#
# simple function to alias with : eal $(make alias)
#
define alias-machine
      function oneview-volume { \
          vol_name=oneview-machine-data; \
          docker info 2> /dev/null|grep "Server Version:.*swarm.*" >/dev/null 2<&1 && IS_SWARM=1 || IS_SWARM=0; \
          if ! docker volume ls -q| grep "^\(.*[\/]\)\?$${vol_name}$$" > /dev/null 2<&1; then \
              echo "Creating volume $${vol_name}"; \
              docker volume create --name=$${vol_name}; \
              docker volume ls -q | \
                  grep "^\(.*[\/]\)\?$${vol_name}$$" | \
                 shuf -n $$(($$(docker volume ls -q |grep "^\(.*[\/]\)\?$${vol_name}$$" | wc -l)-1)) | \
                  xargs -i docker volume rm {}; \
              if [ "$$IS_SWARM" = "1" ]; then \
                  node=$$(docker volume ls|grep ".*${vol_name}$$"|awk "{print $$2}"|awk -F"/" "{print $$1}" | head -1 ); \
                 export DINDC_VOLUME_CONSTRAINT="constraint:node==$${node}"; \
              fi; \
          fi; \
    }; \
    function oneview-machine { \
        envfile=$$(mktemp); \
        oneview-volume; \
        docker_env=$$(bash -c $(PREFIX)/env-opts.sh $$envfile); \
        args=$$@; \
        entrypoint="/usr/local/bin/docker-machine"; \
        if [ -z "$$args" ]; then entrypoint="bash"; fi; \
        eval docker run -it \
                       --rm \
                       -v oneview-machine-data:"$$HOME/.docker" \
                       "$$docker_env" \
                       --entrypoint $$entrypoint \
                       oneview-machine $$debug_opt $$args; \
        rm -f $$envfile; \
    };
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

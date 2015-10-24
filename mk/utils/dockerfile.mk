# goal for this make file is to generate dockerfile from an upstream project

DOCKER_FILE ?=
DOCKER_FILE_URL ?=

include mk/utils/proxy.mk

get-upstream-dockerfile:
		# get the dockerfile.machine from github.com/docker/machine
		curl -s $(DOCKER_FILE_URL) > $(DOCKER_FILE)

gen-dockerfile: proxy-config get-upstream-dockerfile
		echo 'generating docker file'

# goal for this make file is to generate dockerfile from an upstream project

DOCKER_FILE ?=
DOCKER_FILE_URL ?=

PROXY_CONFIG_CONTENT ?= $(shell cat $(PROXY_DOCKER_ENV_FILE))

include mk/utils/proxy.mk

get-upstream-dockerfile:
		# get the dockerfile.machine from github.com/docker/machine
		curl -s $(DOCKER_FILE_URL) > $(DOCKER_FILE)

gen-dockerfile: proxy-config get-upstream-dockerfile
		echo 'generating docker file'
		echo 'cat $(PROXY_DOCKER_ENV_FILE)'
		sed "/FROM.*/ r $(PROXY_DOCKER_ENV_FILE)" $(DOCKER_FILE) > $(DOCKER_FILE).t && mv $(DOCKER_FILE).t $(DOCKER_FILE)
		# setup workdir and add current folder as /go/src/github.com/$GH_USER/$GH_REPO
		sed "s#WORKDIR.*#WORKDIR /go/src/github.com/$(GH_USER)/$(GH_REPO)#g" $(DOCKER_FILE) > $(DOCKER_FILE).t && mv $(DOCKER_FILE).t $(DOCKER_FILE)
		sed "s#ADD.*#ADD . /go/src/github.com/$(GH_USER)/$(GH_REPO)#g" $(DOCKER_FILE) > $(DOCKER_FILE).t && mv $(DOCKER_FILE).t $(DOCKER_FILE)

# # Plain make targets if not requested inside a container
PROXY_CONFIG := ENV HTTP_PROXY $(HTTP_PROXY)\\nENV HTTPS_PROXY $(HTTPS_PROXY)\\nENV http_proxy $(http_proxy)\\nENV https_proxy $(https_proxy)\\nENV no_proxy $(no_proxy)

godeps:
		@echo "Pulling required packages"
		go get github.com/docker/machine
		#TODO : change this once we release and remove the key
		mkdir -p $(PREFIX)/vendor/go/src
		git clone git@github.com:HewlettPackard/oneview-golang.git $(PREFIX)/go/src

include Makefile.inc
ifneq (,$(findstring test-integration,$(MAKECMDGOALS)))
	include mk/main.mk
else ifeq ($(USE_CONTAINER),)
	include mk/main.mk
else
# Otherwise, with docker, swallow all targets and forward into a container
DOCKER_IMAGE_NAME := "docker-machine-build"
DOCKER_CONTAINER_NAME := docker-machine-build-container
# get the dockerfile from docker/machine project so we stay in sync with the versions they use for go
DOCKER_FILE_URL := "https://raw.githubusercontent.com/docker/machine/master/Dockerfile"
DOCKER_FILE := .dockerfile.machine


build:
test: build
%:
		echo $(HTTP_PROXY)
		echo $(https_proxy)
		# get the dockerfile.machine from github.com/docker/machine
		curl -s $(DOCKER_FILE_URL) > $(DOCKER_FILE)
		# setup proxy values
		sed -i "s#FROM golang:1.5.1#FROM golang:1.5.1\\n$(PROXY_CONFIG)#g" $(DOCKER_FILE)
		sed -i "s#\s+ENV#ENV#g" $(DOCKER_FILE)
		# setup workdir and add current folder as /go/src/github.com/$GH_USER/$GH_REPO
		sed -i "s#WORKDIR.*#WORKDIR /go/src/github.com/$(GH_USER)/$(GH_REPO)#g" $(DOCKER_FILE)
		sed -i "s#ADD.*#ADD . /go/src/github.com/$(GH_USER)/$(GH_REPO)#g" $(DOCKER_FILE)
		docker build -f $(DOCKER_FILE) -t $(DOCKER_IMAGE_NAME) .

		test -z '$(shell docker ps -a | grep $(DOCKER_CONTAINER_NAME))' || docker rm -f $(DOCKER_CONTAINER_NAME)

		docker run --name $(DOCKER_CONTAINER_NAME) \
		    -e DEBUG \
		    -e STATIC \
		    -e VERBOSE \
		    -e BUILDTAGS \
		    -e PARALLEL \
		    -e COVERAGE_DIR \
		    -e TARGET_OS \
		    -e TARGET_ARCH \
		    -e PREFIX \
		    $(DOCKER_IMAGE_NAME) \
		    make $@

		test ! -d bin || rm -Rf bin
		test -z "$(findstring build,$(patsubst cross,build,$@))" || docker cp $(DOCKER_CONTAINER_NAME):/go/src/github.com/docker/machine/bin bin

endif

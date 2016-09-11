# # Plain make targets if not requested inside a container

# make builds always run in containers
USE_CONTAINER ?= true

define noop_targets
	@make -pn | sed -rn '/^[^# \t\.%].*:[^=]?/p'|grep -v '='| grep -v '(%)'| grep -v '/'| awk -F':' '{print $$1}'|sort -u;
endef

# setup ONEVIEW_FROM_SRC=1 to use oneview-golang library
# from local path ../oneview-golang instead of installing with go get
# only works on USE_CONTAINER=true
ifeq ($(ONEVIEW_FROM_SRC),1)
	SRC_LIB_ONEVIEW_PATH=$(PREFIX)/../oneview-golang
	DST_LIB_ONEVIEW_PATH=/src/oneview-golang
	DOCKER_VOLUME_ONEVIEW_OPT := --volumes-from oneview-golang
	GO_LIB_ONEVIEW = $(DST_LIB_ONEVIEW_PATH):github.com/HewlettPackard/oneview-golang
endif

#
# supports a way to copy from source the oneview-golang
# lib into the current container project.  Useful for coding on oneview-golang
define golib-setup-from-src
    if [ "$(shell echo $(GO_LIB_ONEVIEW) | cut -c1-1)" = "/" ] ; then \
        echo "skiping go get ..."; \
        echo "install oneview-golang from source --> $(SRC_LIB_ONEVIEW_PATH)"; \
        test -z "$(shell docker ps -a --format '{{.Names}}' | grep '$(DOCKER_CONTAINER_NAME)$$')" || \
            docker rm -f $(DOCKER_CONTAINER_NAME); \
        test -z "$(shell docker ps -a --format '{{.Names}}' | grep 'oneview-golang$$')" || \
            docker rm -f oneview-golang && \
            test -z "$(shell docker volume ls | awk '{print $2}' | grep 'oneview-golang-vol$$')" || \
                docker volume rm oneview-golang-vol; \
        docker volume create --name oneview-golang-vol && \
        docker run -d --name oneview-golang -v oneview-golang-vol:$(DST_LIB_ONEVIEW_PATH) alpine top && \
        cd "$(SRC_LIB_ONEVIEW_PATH)" && \
        docker cp . oneview-golang:$(DST_LIB_ONEVIEW_PATH) && \
        docker kill oneview-golang && \
        echo 'completed creating volume-from for $(DST_LIB_ONEVIEW_PATH)'; \
    else \
        echo "install oneview-golang with go get"; \
    fi;
endef


include Makefile.inc

ifneq (,$(findstring test-integration,$(MAKECMDGOALS)))
	include mk/main.mk
else ifneq (,$(findstring release,$(MAKECMDGOALS)))
	include mk/main.mk
else ifeq ($(USE_CONTAINER),false)
	include mk/main.mk
else

# Otherwise, with docker, swallow all targets and forward into a container
DOCKER_IMAGE_NAME := "docker-machine-build"
DOCKER_CONTAINER_NAME := docker-machine-build-container
# get the dockerfile from docker/machine project so we stay in sync with the versions they use for go
# TODO: delete DOCKER_FILE_URL := "https://raw.githubusercontent.com/docker/machine/master/Dockerfile"
DOCKER_FILE_URL := file://$(PREFIX)/Dockerfile
DOCKER_FILE := .dockerfile.machine

noop:
	@echo When using 'USE_CONTAINER' use a "make <target>"
	@echo
	@echo Possible targets
	@echo
	$(call noop_targets)

oneview-golang-src:
	$(golib-setup-from-src)

clean: gen-dockerfile
build: gen-dockerfile golang-in-docker-build oneview-golang-src
test: gen-dockerfile
golang-in-docker-build:
	docker build -f $(DOCKER_FILE) -t $(DOCKER_IMAGE_NAME) .

%:
		export GO15VENDOREXPERIMENT=1

		test -z "$(shell docker ps -a --format '{{.Names}}' | grep '$(DOCKER_CONTAINER_NAME)$$')" || \
			docker rm -f $(DOCKER_CONTAINER_NAME); \
		docker run --name $(DOCKER_CONTAINER_NAME) \
				$(DOCKER_VOLUME_ONEVIEW_OPT) \
				-e DEBUG \
				-e STATIC \
				-e VERBOSE \
				-e BUILDTAGS \
				-e PARALLEL \
				-e COVERAGE_DIR \
				-e TARGET_OS \
				-e TARGET_ARCH \
				-e PREFIX \
				-e GO15VENDOREXPERIMENT \
				-e TEST_RUN \
				-e ONEVIEW_DEBUG \
				-e GH_USER \
				-e GH_REPO \
				-e VERSION \
				-e GO_LIB_ONEVIEW \
				-e ONEVIEW_FROM_SRC \
				-e GITHUB_TOKEN \
				-e USE_CONTAINER=false \
				$(DOCKER_IMAGE_NAME) \
				make $@

		test ! -d bin || rm -Rf bin
		test -z "$(findstring build,$(patsubst cross,build,$@))" || docker cp $(DOCKER_CONTAINER_NAME):/go/src/github.com/$(GH_USER)/$(GH_REPO)/bin bin

endif

include mk/utils/glide.mk
include mk/utils/dockerfile.mk
include mk/tools.mk

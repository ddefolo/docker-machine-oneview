# use PROXY_CONFIG

define DOCKER_HTTPS_PROXY
ENV HTTPS_PROXY %s\n
endef
define DOCKER_HTTP_PROXY
ENV HTTP_PROXY %s\n
endef
define DOCKER_https_proxy
ENV https_proxy %s\n
endef
define DOCKER_http_proxy
ENV http_proxy %s\n
endef
define DOCKER_NO_PROXY
ENV NO_PROXY %s\n
endef
define DOCKER_no_proxy
ENV no_proxy %s\n
endef

ifndef $(HTTPS_PROXY)
DOCKER_PROXY_CONFIG := $(DOCKER_HTTPS_PROXY)
PROXY_CONFIG := '$(HTTPS_PROXY)'
endif
ifndef $(HTTP_PROXY)
DOCKER_PROXY_CONFIG := $(DOCKER_PROXY_CONFIG)$(DOCKER_HTTP_PROXY)
PROXY_CONFIG := $(PROXY_CONFIG) '$(HTTP_PROXY)'
endif
ifndef $(https_proxy)
DOCKER_PROXY_CONFIG := $(DOCKER_PROXY_CONFIG)$(DOCKER_https_proxy)
PROXY_CONFIG := $(PROXY_CONFIG) '$(https_proxy)'
endif
ifndef $(http_proxy)
DOCKER_PROXY_CONFIG := $(DOCKER_PROXY_CONFIG)$(DOCKER_http_proxy)
PROXY_CONFIG := $(PROXY_CONFIG) '$(http_proxy)'
endif
ifndef $(NO_PROXY)
DOCKER_PROXY_CONFIG := $(DOCKER_PROXY_CONFIG)$(DOCKER_NO_PROXY)
PROXY_CONFIG := $(PROXY_CONFIG) '$(NO_PROXY)'
endif
ifndef $(no_proxy)
DOCKER_PROXY_CONFIG := $(DOCKER_PROXY_CONFIG)$(DOCKER_no_proxy)
PROXY_CONFIG := $(PROXY_CONFIG) '$(no_proxy)'
endif

proxy-clean:
		rm -f .proxy.docker.env

ifneq ($(PROXY_CONFIG),)
proxy-config: proxy-clean
		# generate a
		printf '$(DOCKER_PROXY_CONFIG)' $(PROXY_CONFIG) > .proxy.docker.env
else
proxy-config: proxy-clean
		rm -f .proxy.docker.env
		touch .proxy.docker.env
endif

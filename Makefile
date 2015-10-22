.PHONY: validate

all: validate test build

test:
	script/test

validate-dco:
	script/validate-dco

validate-gofmt:
	script/validate-gofmt

validate: validate-dco validate-gofmt

build:
	script/build

build_local_linux:
	SKIP_USEDOCKER=1 SKIP_BUILD=1 bash script/build -os="linux" -arch="amd64"

PROJECT = "tidy"
PATH := $(PATH)
SHELL := /bin/bash
BUILD_DIR = "build"

BUILDTAGS=debug

default: all

all: build

deps: 
	mkdir -p $(BUILD_DIR)/bin
	go get -tags '$(BUILDTAGS)' -d -v ./...

build: deps; \
        CGO_ENABLED=0 go install -tags '$(BUILDTAGS)' . ; \
        cp $${GOPATH}/bin/tidy $(BUILD_DIR)/bin; \
	cp -r keys $(BUILD_DIR)/bin/;

build-docker: build; \
	TARGET=$(PROJECT):`date +'%Y-%m-%d'`; \
	docker build -t $${TARGET} .

update-key: build; \
	(cd $(BUILD_DIR)/bin/keys/ && ./key-gen.sh)

release: BUILDTAGS=release
release: build

release-docker: BUILDTAGS=release
release-docker: build

update:
	git pull

clean:
	rm -rf build
	go clean -i -r ./...
	git checkout -- .

.PHONY: build clean
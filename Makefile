.PHONY: all \
				check-tool-entr check-tool-go check-tool-sha256sum check-tool-sha1sum check-tool-docker \
				get-verison get-current-sha \
				build build-watch run \
				registry-login image image-publish image-release release-prep release-publish

PROJECT_NAME ?= Neo-Go

VERSION_FILE_PATH ?= .version
VERSION ?= $(shell cat $(VERSION_FILE_PATH))

GIT_REMOTE ?= origin

RELEASES_DIR ?= ./
RELEASE_BINARY_NAME ?= Neo
RELEASE_BUILT_BIN_PATH = ./$(RELEASE_BINARY_NAME)

BIN_NAME = Neo
GO ?= go
ENTR ?= entr
DOCKER ?= docker
SHA256SUM ?= sha256sum
SHA1SUM ?= sha1sum
GIT ?= git

CURRENT_SHA ?= $(shell $(GIT) rev-parse --short HEAD)

all: build

check-tool-entr:
ifeq (, $(shell which $(ENTR)))
	$(error "`entr` is not available please install entr (https://eradman.com/entrproject/)")
endif

check-tool-go:
ifeq (, $(shell which $(GO)))
	$(error "`go` is not available please install go (https://golang.org)")
endif

check-tool-sha256sum:
ifeq (, $(shell which $(SHA256SUM)))
	$(error "`sha256sum` is not available please install sha256sum")
endif

check-tool-sha1sum:
ifeq (, $(shell which $(SHA1SUM)))
	$(error "`sha1sum` is not available please install sha1sum")
endif

check-tool-docker:
ifeq (, $(shell which $(DOCKER)))
	$(error "`docker` is not available please install docker (https://docker.com)")
endif

get-version:
	@echo "$(VERSION)"

get-current-sha:
	@echo "$(CURRENT_SHA)"

build: check-tool-go
	$(GO) build cmd/Neo.Go

build-watch:
	find . -name *.go | $(ENTR) -rc $(GO) build cmd/Neo.Go

build-release: check-tool-go
	CGO_ENABLED=0 $(GO) build cmd/Neo.Go

release-prep:
	$(SHA256SUM) $(RELEASES_DIR)/${RELEASE_BINARY_NAME} > $(RELEASES_DIR)/${RELEASE_BINARY_NAME}.sha256sum
	$(SHA1SUM) $(RELEASES_DIR)/${RELEASE_BINARY_NAME} > $(RELEASES_DIR)/${RELEASE_BINARY_NAME}.sha1sum
	$(GIT) add .
	$(GIT) commit -am "Release $(VERSION)"
	$(GIT) tag $(VERSION) HEAD

release-publish:
	$(GIT) push $(GIT_REMOTE) $(VERSION)

run: build
	./Neo
#############
# Packaging #
#############

REGISTRY_PATH = registry.gitlab.com/visheshc14/$(PROJECT_NAME)

IMAGE_NAME ?= cli

IMAGE_FULL_NAME=${REGISTRY_PATH}/${IMAGE_NAME}:${VERSION}
IMAGE_FULL_NAME_SHA=${REGISTRY_PATH}/${IMAGE_NAME}:${VERSION}-${CURRENT_SHA}

registry-login:
		cat infra/secrets/ci-deploy-token-password.secret | \
		docker login -u $(shell cat infra/secrets/ci-deploy-token-username.secret) --password-stdin registry.gitlab.com

image:
		docker build -f infra/docker/Dockerfile -t $(IMAGE_FULL_NAME_SHA) .

image-publish: check-tool-docker
	$(DOCKER) push ${IMAGE_FULL_NAME_SHA}

image-release:
	$(DOCKER) tag $(IMAGE_FULL_NAME_SHA) $(IMAGE_FULL_NAME)
	$(DOCKER) push $(IMAGE_FULL_NAME)

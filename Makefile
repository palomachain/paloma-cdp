#!/usr/bin/make -f

.PHONY: install build-% docker-% graphql environment-up environment-down test

BINARY_PREFIX=cdp-
BUILD_DIR=./bin
TARGET_OS=linux
TARGET_ARCH=amd64
ENV=CGO_ENABLED=0
VERSION=$(shell git describe --tags --always --dirty)

all: build-ingest build-migrate build-purge build-rest build-transform

install:
	@go mod tidy

build-%:
	@echo "Building $(BINARY_PREFIX)$*"
	@$(ENV) GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build \
		-ldflags "-s -w -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_PREFIX)$* \
		./cmd/$(BINARY_PREFIX)$*

# TODO: Move binary build stage into docker
docker-%: build-%
	@echo "Building docker image for $(BINARY_PREFIX)$*"
	@DOCKER_BUILDKIT=1 docker build \
		--target artifact \
		--build-arg BINARY=$(BINARY_PREFIX)$* \
		-t palomachain/$(BINARY_PREFIX)$*:local \
		-t palomachain/$(BINARY_PREFIX)$*:$(VERSION) \
		-f build/package/Dockerfile \
		.

graphql:
	@go generate ./...

environment-up:
	docker compose -f test/build/docker-compose.yml up -d

environment-down:
	docker compose -f test/build/docker-compose.yml down --volumes

test:
	@gotestsum test ./...

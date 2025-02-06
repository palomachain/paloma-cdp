#!/usr/bin/make -f

.PHONY: install build run graphql environment-up environment-down test

BINARY_NAME=svccdp
BUILD_DIR=./bin
TARGET_OS=linux
TARGET_ARCH=amd64
ENV=CGO_ENABLED=0

install:
	@go mod tidy

build:
	@$(ENV) GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server 

run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

graphql:
	@go generate ./...

environment-up:
	docker compose -f test/build/docker-compose.yml up -d

environment-down:
	docker compose -f test/build/docker-compose.yml down --volumes

test:
	@go test ./...

.DEFAULT_GOAL := build-windows

# globals
BINARY_NAME?=codebeatcli
BUILD_DIR?="./build"
CGO_ENABLED?=0
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date '+%Y-%m-%dT%H:%M:%S %Z')
REPO=github.com/result17/codeBeatCli
VERSION?=<local-build>

# ld flags for go build
LD_FLAGS=-s -w -X '${REPO}/internal/version.BuildDate=${DATE}' -X ${REPO}/internal/version.Commit=${COMMIT} -X ${REPO}/internal/version.Version=${VERSION}

# basic Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# get GOPATH, GOOS and GOARCH according to OS
ifeq ($(OS),Windows_NT) # is Windows_NT on XP, 2000, 7, Vista, 10...
    GOPATH=$(go env GOPATH)
	GOOS := $(shell cmd /c go env GOOS)
	GOARCH := $(shell cmd /c go env GOARCH)
else
    GOPATH=$(shell go env GOPATH)
	GOOS=$(shell go env GOOS)
	GOARCH=$(shell go env GOARCH)
endif

build-all: build-windows
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(MAKE) build-windows

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -v \
		-ldflags "${LD_FLAGS} -X ${REPO}/internal/version.OS=$(GOOS) -X ${REPO}/internal/version.Arch=$(GOARCH)" \
		-o ${BUILD_DIR}/$(BINARY_NAME)-$(GOOS)-$(GOARCH) ./cmd/app

.PHONY: build-windows
build-windows:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -v \
		-ldflags "${LD_FLAGS} -X ${REPO}/internal/version.OS=$(GOOS) -X ${REPO}/internal/version.Arch=$(GOARCH)" \
		-o ${BUILD_DIR}/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe ./cmd/app

install: install-go-modules

.PHONY: install-go-modules
install-go-modules:
	go mod vendor
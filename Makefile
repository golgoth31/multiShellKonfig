#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------

APP_VERSION := latest
PACKAGE ?= $(shell go list ./... | grep "internal/config")
VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_BY=$(shell id -u -n)

GOLINTER:=$(shell command -v golangci-lint 2> /dev/null)

override LDFLAGS += \
  -X ${PACKAGE}.Version=${VERSION} \
  -X ${PACKAGE}.Date=${BUILD_DATE} \
  -X ${PACKAGE}.BuiltBy=${BUILD_BY} \
  -X ${PACKAGE}.GitCommit=${GIT_COMMIT} \


#-----------------------------------------------------------------------------
# BUILD
#-----------------------------------------------------------------------------

.PHONY: default build test publish build_local lint artifact_linux artifact_darwin deploy
default:  test lint build swagger

test:
	go test -v ./...

build_container_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '${LDFLAGS}' -a -o msk-bin main.go

build_container_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '${LDFLAGS}' -a -o msk-bin main.go

build_container_darwin_arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags '${LDFLAGS}' -a -o msk-bin main.go

lint:
	golangci-lint run

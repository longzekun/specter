SERVER_BIN := specter-server
CLIENT_BIN := specter-client

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_DIRTY := $(shell git diff --quiet || echo true)
BUILD_TIME := $(shell date +%s)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

VERSION := 0.1.0

SERVER_LDFLAGS = -ldflags "-s -w \
	-X version.Version=$(VERSION) \
	-X version.Commit=$(GIT_COMMIT) \
	-X version.Dirty=$(GIT_DIRTY) \
	-X version.OS=$(GOOS) \
	-X version.Arch=$(GOARCH)"

build-server:
	go build $(SERVER_LDFLAGS) -o $(SERVER_BIN) server/main.go

build-client:
	go build -o $(CLIENT_BIN) client/main.go

build: build-server build-client

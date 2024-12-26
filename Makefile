BIN_DIR=$(PWD)/bin
PICO_DIR=$(PWD)/cmd/pico-cli
CC=gcc
CXX=g++
VERSION=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "0.0.0")
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags="-s -w -X github.com/PicoTools/pico-cli/internal/version.gitCommit=${BUILD} -X github.com/PicoTools/pico-cli/internal/version.gitVersion=${VERSION}"

.PHONY: pico-cli
pico-cli:
	@mkdir -p ${BIN_DIR}
	@echo "Building operator cli..."
	CGO_ENABLED=0 CC=${CC} CXX=${CXX} go build ${LDFLAGS} -o ${BIN_DIR}/pico-cli ${PICO_DIR}
	@strip bin/pico-cli

.PHONY: go-sync
go-sync:
	@go mod tidy && go mod vendor

.PHONY: dep-shared
dep-shared:
	@echo "Update shared components..."
	@export GOPRIVATE="github.com/PicoTools" && go get -u github.com/PicoTools/pico-shared/ && go mod tidy && go mod vendor

.PHONY: dep-plan
dep-plan:
	@echo "Update mlan components..."
	@export GOPRIVATE="github.com/PicoTools" && go get -u github.com/PicoTools/plan/ && go mod tidy && go mod vendor

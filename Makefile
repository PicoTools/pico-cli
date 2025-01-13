BIN_DIR=$(PWD)/bin
PICO_DIR=$(PWD)/cmd/pico-cli
CC=gcc
CXX=g++
GOFILES=`go list ./...`
GOFILESNOTEST=`go list ./... | grep -v test`
VERSION=$(shell git describe --abbrev=0 --tags 2>/dev/null || echo "v0.0.0")
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags="-s -w -X github.com/PicoTools/pico-cli/internal/version.gitCommit=${BUILD} -X github.com/PicoTools/pico-cli/internal/version.gitVersion=${VERSION}"

build-all: darwin-arm64 darwin-amd64 linux-arm64 linux-amd64

darwin-arm64: go-lint
	@mkdir -p ${BIN_DIR}
	@echo "Building operator cli darwin/arm64 ${VERSION}" 
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 CC=${CC} CXX=${CXX} go build -trimpath ${LDFLAGS} -o ${BIN_DIR}/pico-cli.darwin.arm64 ${PICO_DIR}

darwin-amd64: go-lint
	@mkdir -p ${BIN_DIR}
	@echo "Building operator cli darwin/amd64 ${VERSION}"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 CC=${CC} CXX=${CXX} go build -trimpath ${LDFLAGS} -o ${BIN_DIR}/pico-cli.darwin.amd64 ${PICO_DIR}

linux-arm64: go-lint
	@mkdir -p ${BIN_DIR}
	@echo "Building operator cli linux/arm64 ${VERSION}"
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 CC=${CC} CXX=${CXX} go build -trimpath ${LDFLAGS} -o ${BIN_DIR}/pico-cli.linux.arm64 ${PICO_DIR}

linux-amd64: go-lint
	@mkdir -p ${BIN_DIR}
	@echo "Building operator cli linux/amd64 ${VERSION}"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 CC=${CC} CXX=${CXX} go build -trimpath ${LDFLAGS} -o ${BIN_DIR}/pico-cli.linux.amd64 ${PICO_DIR}

go-lint:
	@echo "Linting Golang code..."
	@go fmt ${GOFILES}
	@go vet ${GOFILESNOTEST}

go-sync:
	@go mod tidy && go mod vendor

dep-shared:
	@echo "Update shared components..."
	@export GOPRIVATE="github.com/PicoTools" && go get -u github.com/PicoTools/pico-shared/ && go mod tidy && go mod vendor

dep-plan:
	@echo "Update plan components..."
	@export GOPRIVATE="github.com/PicoTools" && go get -u github.com/PicoTools/plan/ && go mod tidy && go mod vendor

BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --tags --dirty --exact-match 2>/dev/null || git rev-parse --short HEAD)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X github.com/ci4rail/device-state-service/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=arm64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/device-state-service main.go

clean:
	rm -rf ${BIN_DIR}/device-state-service

.PHONY: clean

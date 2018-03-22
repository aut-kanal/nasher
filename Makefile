## Project Vars ##########################################################
GO_VARS = ENABLE_CGO=0 GOOS=darwin GOARCH=amd64
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := gitlab.com/kanalbot/nasher
ROOT_DIRECTORY := $(shell pwd)
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).Title=nasherd

.PHONY: help clean

nasherd: *.go */*.go */*/*.go Gopkg.lock
	$(GO_VARS) $(GO) build -i -o="nasherd" -ldflags="$(LD_FLAGS)" $(ROOT)/cmd/nasher

clean:
	rm -rf nasherd

help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  nasherd     to build the main binary for current platform"
	@echo "  clean             to remove generated files"

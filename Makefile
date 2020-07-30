NAME := errwrp
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/errwrp

.PHONY: setup dep test main clean install

all: vendor build

vendor:
	go mod vendor

build:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o bin/errwrp cmd/mustwrap/main.go

test:
	go test -v ./...

test/cover:
	go test -v -coverprofile=out ./...

## remove build files
clean:
	rm -rf ./bin/*

NAME := errwrp
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/errwrp

.PHONY: setup dep test main clean install

all: vendor build

vendor:
	go mod vendor

testdata/src/github.com:
	GOPATH=$(CURDIR)/testdata GO111MODULE=off go get github.com/pkg/errors

build:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o bin/errwrp cmd/mustwrap/main.go

test: testdata/src/github.com
	go test -v ./...

test/cover: testdata/src/github.com
	go test -v -coverprofile=out ./...

## remove build files
clean:
	rm -rf ./bin/*

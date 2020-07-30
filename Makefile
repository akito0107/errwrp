NAME := errwrp
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/akito0107/errwrp

.PHONY: setup dep test main clean install testdata

all: vendor build

vendor: go.mod go.sum
	go mod vendor

testdata: testdata/src/github.com testdata/src/golang.org

testdata/src/github.com:
	GOPATH=$(CURDIR)/testdata GO111MODULE=off go get github.com/pkg/errors

testdata/src/golang.org:
	GOPATH=$(CURDIR)/testdata GO111MODULE=off go get golang.org/x/xerrors

build:
	go build -mod=vendor -ldflags "$(LDFLAGS)" -o bin/errwrp cmd/mustwrap/main.go

test: testdata
	go test -v ./...

test/cover: testdata
	go test -v -coverprofile=out ./...

## remove build files
clean:
	rm -rf ./bin/*

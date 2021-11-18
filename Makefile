COMMIT=$(shell git rev-parse HEAD)
#VERSION=$(shell git describe --tags --exact-match --always)
VERSION=v0.0.0
DATE=$(shell date +'%FT%TZ%z')

dist/html0.so:  $(shell find . -type f -name '*.go')
	go build -buildmode=c-shared -tags="shared" \
	-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
	-o $@  shared.go


all: dist/html0.so

test:
	./test.sh 

format:
	gofmt -s -w .

.PHONY: httpbin all test format
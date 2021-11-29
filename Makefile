COMMIT=$(shell git rev-parse HEAD)
#VERSION=$(shell git describe --tags --exact-match --always)
VERSION=v0.0.0
DATE=$(shell date +'%FT%TZ%z')

dist/html0-macos.dylib:  $(shell find . -type f -name '*.go')
	go build -buildmode=c-shared -tags="shared" \
	-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
	-o $@ shared.go

dist/html0-linux.so:  $(shell find . -type f -name '*.go')
	go build -buildmode=c-shared -tags="shared" \
	-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
	-o $@ shared.go

dist/html0-windows.dll:  $(shell find . -type f -name '*.go')
	go build -buildmode=c-shared -tags="shared" \
	-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
	-o $@ shared.go

macos: dist/html0-macos.dylib

linux: dist/html0-linux.so

windows: dist/html0-windows.dll

test:
	 python3 test.py

format:
	gofmt -s -w .

.PHONY: test format linux macos windows 
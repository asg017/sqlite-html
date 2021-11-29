COMMIT=$(shell git rev-parse HEAD)
#VERSION=$(shell git describe --tags --exact-match --always)
VERSION=v0.0.0
DATE=$(shell date +'%FT%TZ%z')


# from riyaz-ali/sqlite
export CGO_LDFLAGS = -Wl,--unresolved-symbols=ignore-in-object-files
ifeq ($(shell uname -s),Darwin)
	export CGO_LDFLAGS = -Wl,-undefined,dynamic_lookup
endif

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

#dist/html0-linux-arm.so: 
#	go build -buildmode=c-shared -tags="shared" \
#		-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' \
#		-o $@ shared.go
#
#linux-docker-arm: 
#	docker run \
#		--rm -it \
#		-v $PWD:/mnt -w /mnt \
#		-e CGO_ENABLED=1 -e CC=arm-linux-gnueabi-gcc \
#		-e GOOS=linux -e GOARCH=arm \
#		golang:1.17-buster \
#		make dist/html0-linux-arm.so


#
#linux: dist/html0-linux.so
#
#linux-docker: 
#	docker run --rm -it -v $PWD:/mnt -w /mnt golang:1.17-buster make dist/html0-linux.so
#
#dist-linux-arm:
#	docker run --rm \
#		-v "${PWD}":/usr/src/litestream \
#		-w /usr/src/litestream \
#		-e CGO_ENABLED=1 -e CC=arm-linux-gnueabihf-gcc -e GOOS=linux -e GOARCH=arm \
#		golang-xc:1.16 \
#		go build -v -o dist/litestream-linux-arm ./cmd/litestream

#dist-linux-arm64:
#	docker run --rm \
#		-v "${PWD}":/usr/src/litestream \
#		-w /usr/src/litestream \
#		-e CGO_ENABLED=1 -e CC=aarch64-linux-gnu-gcc -e GOOS=linux -e GOARCH=arm64 \
#		golang-xc:1.16 \
#		x`go build -v -o dist/litestream-linux-arm64 ./cmd/litestream

#all: dist/html0-macos.so

test:
	 python3 test.py

format:
	gofmt -s -w .

.PHONY: all test format linux mac windows #linux-docker-arm
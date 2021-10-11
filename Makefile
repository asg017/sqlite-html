dist/dom.so:  $(shell find . -type f -name '*.go')
	go build -buildmode=c-shared -o $@ -tags="shared" shared.go
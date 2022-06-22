COMMIT=$(shell git rev-parse HEAD)
#VERSION=$(shell git describe --tags --exact-match --always)
VERSION=v0.0.0
DATE=$(shell date +'%FT%TZ%z')

GO_BUILD_LDFLAGS=-ldflags '-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)' 
GO_BUILD_CGO_CFLAGS=CGO_CFLAGS=-DSQLITE3_INIT_FN=sqlite3_html_init

ifeq ($(OS),Windows_NT)
CONFIG_WINDOWS=y
endif

ifeq ($(shell uname -s),Darwin)
CONFIG_DARWIN=y
else
CONFIG_LINUX=y
endif

# framework stuff is needed bc https://github.com/golang/go/issues/42459#issuecomment-896089738                                                                           
ifdef CONFIG_DARWIN
LOADABLE_EXTENSION=dylib
SQLITE3_CFLAGS=-framework CoreFoundation -framework Security
endif

ifdef CONFIG_LINUX
LOADABLE_EXTENSION=so
endif


ifdef CONFIG_WINDOWS
LOADABLE_EXTENSION=dll
endif

TARGET_LOADABLE=dist/html0.$(LOADABLE_EXTENSION)
TARGET_OBJ=dist/html0.o
TARGET_SQLITE3=dist/sqlite3
TARGET_PACKAGE=dist/package.zip

loadable: $(TARGET_LOADABLE)
sqlite3: $(TARGET_SQLITE3)
package: $(TARGET_PACKAGE)
all: loadable sqlite3 package

$(TARGET_LOADABLE):  $(shell find . -type f -name '*.go')
	$(GO_BUILD_CGO_CFLAGS) go build -buildmode=c-shared -tags="shared" \
	$(GO_BUILD_LDFLAGS) \
	-o $@ .

$(TARGET_OBJ):  $(shell find . -type f -name '*.go')
	$(GO_BUILD_CGO_CFLAGS) CGO_ENABLED=1 go build -buildmode=c-archive \
	$(GO_BUILD_LDFLAGS) \
	-o $@ .

# I don't think we can include DSQLITE_OMIT_LOAD_EXTENSION - maybe riyaz-ali/sqlite uses it?
# add back later -DHAVE_READLINE -lreadline -lncurses
$(TARGET_SQLITE3): $(TARGET_OBJ) dist/sqlite3-extra.c sqlite/shell.c
	gcc \
	$(SQLITE3_CFLAGS) \
	-lm -pthread \
	dist/sqlite3-extra.c sqlite/shell.c $(TARGET_OBJ) \
	-ldl -L. -I./sqlite \
	-DSQLITE_EXTRA_INIT=core_init -DSQLITE3_INIT_FN=sqlite3_html_init \
	-o $@

$(TARGET_PACKAGE): $(TARGET_LOADABLE) $(TARGET_OBJ) sqlite/sqlite-html.h $(TARGET_SQLITE3)
	zip --junk-paths $@ $(TARGET_LOADABLE) $(TARGET_OBJ) sqlite/sqlite-html.h $(TARGET_SQLITE3)

dist/sqlite3-extra.c: sqlite/sqlite3.c sqlite/core_init.c
	cat sqlite/sqlite3.c sqlite/core_init.c > $@

clean:
	rm dist/*

test:
	make test-loadable
	make test-sqlite3

test-loadable:
	python3 tests/test-loadable.py

test-sqlite3:
	python3 tests/test-sqlite3.py

format:
	gofmt -s -w .

.PHONY: all clean format \
	test test-loadable test-sqlite3 \
	loadable sqlite3 package

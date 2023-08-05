COMMIT=$(shell git rev-parse HEAD)
VERSION=$(shell cat VERSION)
DATE=$(shell date +'%FT%TZ%z')
VENDOR_SQLITE=$(shell pwd)/sqlite
GO_BUILD_LDFLAGS=-ldflags '-X main.Version=v$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)'
#GO_BUILD_CGO_CFLAGS=CGO_CFLAGS=-DSQLITE3_INIT_FN=sqlite3_html_init
GO_BUILD_CGO_CFLAGS=CGO_ENABLED=1 CGO_CFLAGS="-DUSE_LIBSQLITE3" CPATH="$(VENDOR_SQLITE)"


ifeq ($(shell uname -s),Darwin)
CONFIG_DARWIN=y
else ifeq ($(OS),Windows_NT)
CONFIG_WINDOWS=y
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

ifdef python
PYTHON=$(python)
else
PYTHON=python3
endif

ifdef IS_MACOS_ARM
RENAME_WHEELS_ARGS=--is-macos-arm
else
RENAME_WHEELS_ARGS=
endif

prefix=dist

TARGET_LOADABLE=$(prefix)/html0.$(LOADABLE_EXTENSION)
TARGET_OBJ=$(prefix)/html0.o
TARGET_WHEELS=$(prefix)/wheels
TARGET_SQLITE3=$(prefix)/sqlite3
TARGET_PACKAGE=$(prefix)/package.zip

INTERMEDIATE_PYPACKAGE_EXTENSION=python/sqlite_html/sqlite_html/html0.$(LOADABLE_EXTENSION)

loadable: $(TARGET_LOADABLE)
sqlite3: $(TARGET_SQLITE3)
package: $(TARGET_PACKAGE)
all: loadable sqlite3 package

$(prefix):
	mkdir -p $(prefix)

$(TARGET_WHEELS): $(prefix)
	mkdir -p $(TARGET_WHEELS)

$(TARGET_LOADABLE):  $(shell find . -type f -name '*.go')
	$(GO_BUILD_CGO_CFLAGS) go build -buildmode=c-shared -tags="shared" \
	$(GO_BUILD_LDFLAGS) \
	-o $@ .

python: $(TARGET_WHEELS) $(TARGET_LOADABLE) $(TARGET_WHEELS) scripts/rename-wheels.py $(shell find python/sqlite_html -type f -name '*.py')
	cp $(TARGET_LOADABLE) $(INTERMEDIATE_PYPACKAGE_EXTENSION)
	rm $(TARGET_WHEELS)/sqlite_html* || true
	pip3 wheel python/sqlite_html/ -w $(TARGET_WHEELS)
	python3 scripts/rename-wheels.py $(TARGET_WHEELS) $(RENAME_WHEELS_ARGS)
	echo "✅ generated python wheel"

python-versions: python/version.py.tmpl
	VERSION=$(VERSION) envsubst < python/version.py.tmpl > python/sqlite_html/sqlite_html/version.py
	echo "✅ generated python/sqlite_html/sqlite_html/version.py"

	VERSION=$(VERSION) envsubst < python/version.py.tmpl > python/datasette_sqlite_html/datasette_sqlite_html/version.py
	echo "✅ generated python/datasette_sqlite_html/datasette_sqlite_html/version.py"

datasette: $(TARGET_WHEELS) $(shell find python/datasette_sqlite_html -type f -name '*.py')
	rm $(TARGET_WHEELS)/datasette* || true
	pip3 wheel python/datasette_sqlite_html/ --no-deps -w $(TARGET_WHEELS)

npm: VERSION npm/platform-package.README.md.tmpl npm/platform-package.package.json.tmpl npm/sqlite-html/package.json.tmpl scripts/npm_generate_platform_packages.sh
	scripts/npm_generate_platform_packages.sh

deno: VERSION deno/deno.json.tmpl
	scripts/deno_generate_package.sh

bindings/ruby/lib/version.rb: bindings/ruby/lib/version.rb.tmpl VERSION
	VERSION=$(VERSION) envsubst < $< > $@

ruby: bindings/ruby/lib/version.rb

version:
	make python-versions
	make python
	make npm
	make deno
	make ruby

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
	-L. -I./sqlite \
	-DSQLITE_THREADSAFE=0 -DSQLITE_OMIT_LOAD_EXTENSION=1 \
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
	make test-python
	make test-npm
	make test-deno

test-loadable: $(TARGET_LOADABLE)
	$(PYTHON) tests/test-loadable.py

test-python:
	$(PYTHON) tests/test-python.py

test-npm:
	node npm/sqlite-html/test.js

test-deno:
	deno task --config deno/deno.json test

test-loadable-watch: $(TARGET_LOADABLE)
	watchexec -w . -w $(TARGET_LOADABLE) -w tests/test-loadable.py --clear -- make test-loadable

test-sqlite3: sqlite3
	python3 tests/test-sqlite3.py

format:
	gofmt -s -w .

publish-release:
	./scripts/publish_release.sh

.PHONY: all clean format publish-release \
	python python-versions datasette npm deno ruby version \
	test test-loadable test-sqlite3 \
	loadable sqlite3 package

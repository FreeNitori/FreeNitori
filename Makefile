all: deps assets build
run: assets build start

LDFLAGS = "-s -w -X 'git.randomchars.net/RandomChars/FreeNitori/nitori/state.version=$(shell echo -n `git describe --tags`; if ! [ "`git status -s`" = '' ]; then echo -n '-dirty'; fi)' -X 'git.randomchars.net/RandomChars/FreeNitori/nitori/state.revision=$(shell git rev-parse --short HEAD)'"
ifeq ($(shell go env GOOS), windows)
   Suffix = ".exe"
endif

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	@echo "Packaging assets..."
	@$$(go env GOPATH)/bin/go-bindata -o binaries/confdefault/confdefault.go -pkg confdefault -prefix assets/ ./assets/nitori.conf
	@$$(go env GOPATH)/bin/go-bindata -o binaries/tmpl/tmpl.go -pkg tmpl -prefix assets/web/templates/ ./assets/web/templates/*
	@$$(go env GOPATH)/bin/go-bindata -fs -o binaries/static/static.go -pkg static -prefix assets/web/public/ ./assets/web/public/...

.PHONY: plugins
plugins:
	@echo "Building plugins..."
	@for pl in $(shell sh -c "ls plugins/*/main.go"); do go build -ldflags="-s -w" --buildmode=plugin -o ./plugins $$PWD/$${pl::-7}; done;

.PHONY: build
build:
	@echo "Building FreeNitori..."
	@go build -tags=jsoniter -ldflags=$(LDFLAGS) -o build/freenitori$(Suffix) $$PWD/cmd/server
	@echo "Building nitorictl..."
	@go build -tags=jsoniter -ldflags=$(LDFLAGS) -o build/nitorictl$(Suffix) $$PWD/cmd/cli

.PHONY: start
start:
	@./build/freenitori

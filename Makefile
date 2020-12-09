all: deps assets build
run: assets build start

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@/usr/bin/env go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	@echo "Packaging assets..."
	@$$(go env GOPATH)/bin/go-bindata -o binaries/confdefault/confdefault.go -pkg confdefault -prefix assets/ ./assets/nitori.conf
	@$$(go env GOPATH)/bin/go-bindata -o binaries/tmpl/tmpl.go -pkg tmpl -prefix assets/web/templates/ ./assets/web/templates/*
	@$$(go env GOPATH)/bin/go-bindata -fs -o binaries/static/static.go -pkg static -prefix assets/web/static/ ./assets/web/static/...

.PHONY: plugins
plugins:
	@echo "Building plugins..."
	@for pl in $(shell sh -c "ls plugins/*/main.go"); do /usr/bin/env go build -ldflags="-s -w" --buildmode=plugin -o ./plugins $$PWD/$${pl::-7}; done;

.PHONY: internal
internal:
	@echo "Building internal plugins..."
	@for pl in $(shell ls "internal/"); do /usr/bin/env go build -ldflags="-s -w" --buildmode=plugin -o ./plugins $$PWD/internal/$$pl; echo "Built $${pl}."; done;

.PHONY: build
build: internal
	@echo "Building FreeNitori..."
	@/usr/bin/env go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori $$PWD/server

.PHONY: start
start:
	@./build/freenitori
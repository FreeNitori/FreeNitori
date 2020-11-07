all: deps assets build
run: build assets start

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@/usr/bin/env go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	@echo "Packaging assets..."
	@$$GOPATH/bin/go-bindata -o proc/supervisor/confdefault/confdefault.go -pkg confdefault -prefix assets/ ./assets/nitori.conf
	@$$GOPATH/bin/go-bindata -o proc/webserver/tmpl/tmpl.go -pkg tmpl -prefix assets/web/templates/ ./assets/web/templates/*
	@$$GOPATH/bin/go-bindata -fs -o proc/webserver/static/static.go -pkg static -prefix assets/web/static/ ./assets/web/static/...

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
	@for proc in $(shell ls "proc/"); do /usr/bin/env go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-$$proc $$PWD/proc/$$proc; echo "Built $${proc}."; done;
	@cp -f build/freenitori-supervisor build/freenitori

.PHONY: start
start:
	@./build/freenitori
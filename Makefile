all: deps build
run: build start

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	@echo "Packaging assets..."
	@$$GOPATH/bin/go-bindata -o nitori/config/assets.go -pkg config -prefix assets/ ./assets/*
	@$$GOPATH/bin/go-bindata -o proc/webserver/tmpl/tmpl.go -pkg tmpl -prefix assets/web/templates/ ./assets/web/templates/*
	@$$GOPATH/bin/go-bindata -fs -o proc/webserver/static/static.go -pkg static -prefix assets/web/static/ ./assets/web/static/...

.PHONY: start
start:
	@./build/freenitori

.PHONY: build
build: assets
	@echo "Building FreeNitori..."
	@for proc in $(shell ls "proc/"); do go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-$$proc $$PWD/proc/$$proc; echo "Built $${proc}."; done;
	@cp -f build/freenitori-supervisor build/freenitori
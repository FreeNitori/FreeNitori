all: build
test: build-test run-test

.PHONY: deps
deps:
	go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets:
	$$GOPATH/bin/go-bindata -o nitori/config/assets.go -pkg config -prefix assets/ ./assets/* ./assets/web/templates/*
	$$GOPATH/bin/go-bindata -fs -o proc/webserver/static/static.go -pkg static -prefix assets/web/static/ ./assets/web/static/...

.PHONY: run-test
run-test:
	./build/freenitori

.PHONY: build
build: deps assets
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-supervisor $$PWD/proc/supervisor
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-chatbackend $$PWD/proc/chatbackend
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-webserver $$PWD/proc/webserver
	go build -tags=jsoniter -ldflags="-s -w" -o build/freenitori-shell $$PWD/proc/shell
	cp -f build/freenitori-supervisor build/freenitori

.PHONY: build-test
build-test: assets
	go build -tags=jsoniter -o build/freenitori-supervisor $$PWD/proc/supervisor
	go build -tags=jsoniter -o build/freenitori-chatbackend $$PWD/proc/chatbackend
	go build -tags=jsoniter -o build/freenitori-webserver $$PWD/proc/webserver
	go build -tags=jsoniter -o build/freenitori-shell $$PWD/proc/shell
	cp -f build/freenitori-supervisor build/freenitori
